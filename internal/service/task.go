package service

import (
	"context"
	"golang-service-template/internal/dao/model"
	"golang-service-template/internal/dao/query"
	"golang-service-template/internal/errz"
	"golang-service-template/internal/telemetry"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type TaskService interface {
	Create(ctx context.Context, task model.Task) (*model.Task, error)
	Get(ctx context.Context, id string) (*model.Task, error)
	Find(ctx context.Context) ([]*model.Task, error)
	FindByUserId(ctx context.Context, userId string) ([]*model.Task, error)
	Update(ctx context.Context, id string, entity map[string]any) (*model.Task, error)
	Delete(ctx context.Context, id string) error
}

type taskService struct {
	db        *gorm.DB
	q         *query.Query
	redis     *redis.Client
	telemetry *telemetry.Telemetry
}

func NewTaskService(i *do.Injector) (TaskService, error) {
	db := do.MustInvoke[*gorm.DB](i)
	tel, _ := do.Invoke[*telemetry.Telemetry](i) // Optional telemetry

	return &taskService{
		db:        db,
		q:         query.Use(db),
		redis:     do.MustInvoke[*redis.Client](i),
		telemetry: tel,
	}, nil
}

// Create implements TaskService.
func (s *taskService) Create(ctx context.Context, entity model.Task) (*model.Task, error) {
	start := time.Now()

	// Create trace span using generic method
	ctx, span := s.telemetry.CreateSpan(ctx, "task_create",
		attribute.String("operation", "create"))
	defer span.End()

	newID, err := uuid.NewV7()
	if err != nil {
		s.telemetry.RecordError(ctx, err)
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to generate new id", err)
	}

	entityp := &entity
	entityp.ID = newID.String()
	entityp.CreatedBy = newID.String() // TODO: get user id from context

	if err := query.Use(s.db).WithContext(ctx).Task.Create(entityp); err != nil {
		// Record error metrics using generic methods
		s.telemetry.Increment(ctx, "task_create_total",
			attribute.String("status", "error"))
		s.telemetry.RecordDuration(ctx, "task_create_duration_seconds",
			start,
			attribute.String("status", "error"))
		s.telemetry.RecordError(ctx, err)
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to create task", err)
	}

	// Record success metrics using generic methods
	s.telemetry.Increment(ctx, "task_create_total",
		attribute.String("status", "success"))
	s.telemetry.RecordDuration(ctx, "task_create_duration_seconds",
		start,
		attribute.String("status", "success"))

	// Add task ID to span now that we have it
	span.SetAttributes(attribute.String("task.id", entityp.ID))

	return entityp, nil
}

// Get implements TaskService.
func (s *taskService) Get(ctx context.Context, id string) (*model.Task, error) {
	start := time.Now()

	// Create trace span using generic method
	ctx, span := s.telemetry.CreateSpan(ctx, "task_get",
		attribute.String("operation", "get"),
		attribute.String("task.id", id))
	defer span.End()

	entity, err := s.q.WithContext(ctx).Task.Where(s.q.Task.ID.Eq(id)).First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.telemetry.Increment(ctx, "task_get_total",
			attribute.String("status", "not_found"))
		s.telemetry.RecordDuration(ctx, "task_get_duration_seconds",
			start,
			attribute.String("status", "not_found"))
		return nil, errz.NewPrettyError(http.StatusNotFound, "not_found", "entity not found", err)
	}

	if err != nil {
		s.telemetry.Increment(ctx, "task_get_total",
			attribute.String("status", "error"))
		s.telemetry.RecordDuration(ctx, "task_get_duration_seconds",
			start,
			attribute.String("status", "error"))
		s.telemetry.RecordError(ctx, err)
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entity", err)
	}

	// Authorization check: verify user owns the task
	if userId := ctx.Value("context_key_user_id"); userId != nil {
		if userIdStr, ok := userId.(string); ok && entity.CreatedBy != userIdStr {
			s.telemetry.Increment(ctx, "task_get_total",
				attribute.String("status", "forbidden"))
			s.telemetry.RecordDuration(ctx, "task_get_duration_seconds",
				start,
				attribute.String("status", "forbidden"))
			return nil, errz.NewPrettyError(http.StatusForbidden, "forbidden", "you don't have permission to access this task", nil)
		}
	}

	// Record success
	s.telemetry.Increment(ctx, "task_get_total",
		attribute.String("status", "success"))
	s.telemetry.RecordDuration(ctx, "task_get_duration_seconds",
		start,
		attribute.String("status", "success"))
	return entity, nil
}

// GetAll implements TaskService.
func (s *taskService) Find(ctx context.Context) ([]*model.Task, error) {
	start := time.Now()

	// Create trace span using generic method
	ctx, span := s.telemetry.CreateSpan(ctx, "task_find_all",
		attribute.String("operation", "find_all"))
	defer span.End()

	entities, err := s.q.WithContext(ctx).Task.Find()

	if err != nil {
		s.telemetry.Increment(ctx, "task_find_all_total",
			attribute.String("status", "error"))
		s.telemetry.RecordDuration(ctx, "task_find_all_duration_seconds",
			start,
			attribute.String("status", "error"))
		s.telemetry.RecordError(ctx, err)
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entities", err)
	}

	// Record success
	s.telemetry.Increment(ctx, "task_find_all_total",
		attribute.String("status", "success"))
	s.telemetry.RecordDuration(ctx, "task_find_all_duration_seconds",
		start,
		attribute.String("status", "success"))
	span.SetAttributes(attribute.Int("task.count", len(entities)))
	return entities, nil
}

// GetAll implements TaskService.
func (s *taskService) FindByUserId(ctx context.Context, userId string) ([]*model.Task, error) {
	start := time.Now()

	// Create trace span using generic method
	ctx, span := s.telemetry.CreateSpan(ctx, "task_find_by_user",
		attribute.String("operation", "find_by_user"),
		attribute.String("user.id", userId))
	defer span.End()

	entities, err := s.q.WithContext(ctx).Task.Where(s.q.Task.CreatedBy.Eq(userId)).Find()

	if err != nil {
		s.telemetry.Increment(ctx, "task_find_by_user_total",
			attribute.String("status", "error"))
		s.telemetry.RecordDuration(ctx, "task_find_by_user_duration_seconds",
			start,
			attribute.String("status", "error"))
		s.telemetry.RecordError(ctx, err)
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entities", err)
	}

	// Record success
	s.telemetry.Increment(ctx, "task_find_by_user_total",
		attribute.String("status", "success"))
	s.telemetry.RecordDuration(ctx, "task_find_by_user_duration_seconds",
		start,
		attribute.String("status", "success"))
	span.SetAttributes(attribute.Int("task.count", len(entities)))
	return entities, nil
}

// Update implements TaskService.
// using map here to avoid headache of handling Go's zero value
// we pass whatever passed validation in handler
func (s *taskService) Update(ctx context.Context, id string, entity map[string]any) (*model.Task, error) {
	start := time.Now()

	// Create trace span using generic method
	ctx, span := s.telemetry.CreateSpan(ctx, "task_update",
		attribute.String("operation", "update"),
		attribute.String("task.id", id))
	defer span.End()

	// Authorization check: verify user owns the task before updating
	if userId := ctx.Value("context_key_user_id"); userId != nil {
		existingTask, err := s.q.WithContext(ctx).Task.Where(s.q.Task.ID.Eq(id)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.telemetry.Increment(ctx, "task_update_total",
					attribute.String("status", "not_found"))
				s.telemetry.RecordDuration(ctx, "task_update_duration_seconds",
					start,
					attribute.String("status", "not_found"))
				return nil, errz.NewPrettyError(http.StatusNotFound, "not_found", "entity not found", err)
			}
			s.telemetry.Increment(ctx, "task_update_total",
				attribute.String("status", "error"))
			s.telemetry.RecordDuration(ctx, "task_update_duration_seconds",
				start,
				attribute.String("status", "error"))
			s.telemetry.RecordError(ctx, err)
			return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to check task ownership", err)
		}
		
		if userIdStr, ok := userId.(string); ok && existingTask.CreatedBy != userIdStr {
			s.telemetry.Increment(ctx, "task_update_total",
				attribute.String("status", "forbidden"))
			s.telemetry.RecordDuration(ctx, "task_update_duration_seconds",
				start,
				attribute.String("status", "forbidden"))
			return nil, errz.NewPrettyError(http.StatusForbidden, "forbidden", "you don't have permission to update this task", nil)
		}
	}

	_, err := s.q.WithContext(ctx).Task.Where(s.q.Task.ID.Eq(id)).Updates(entity)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.telemetry.Increment(ctx, "task_update_total",
			attribute.String("status", "not_found"))
		s.telemetry.RecordDuration(ctx, "task_update_duration_seconds",
			start,
			attribute.String("status", "not_found"))
		return nil, errz.NewPrettyError(http.StatusNotFound, "not_found", "entity not found", err)
	}

	if err != nil {
		s.telemetry.Increment(ctx, "task_update_total",
			attribute.String("status", "error"))
		s.telemetry.RecordDuration(ctx, "task_update_duration_seconds",
			start,
			attribute.String("status", "error"))
		s.telemetry.RecordError(ctx, err)
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to update entities", err)
	}

	// Record success for update operation
	s.telemetry.Increment(ctx, "task_update_total",
		attribute.String("status", "success"))
	s.telemetry.RecordDuration(ctx, "task_update_duration_seconds",
		start,
		attribute.String("status", "success"))

	// Get the updated entity (this will have its own telemetry and authorization check)
	return s.Get(ctx, id)
}

// Delete implements TaskService.
func (s *taskService) Delete(ctx context.Context, id string) error {
	start := time.Now()

	// Create trace span using generic method
	ctx, span := s.telemetry.CreateSpan(ctx, "task_delete",
		attribute.String("operation", "delete"),
		attribute.String("task.id", id))
	defer span.End()

	// Authorization check: verify user owns the task before deleting
	if userId := ctx.Value("context_key_user_id"); userId != nil {
		existingTask, err := s.q.WithContext(ctx).Task.Where(s.q.Task.ID.Eq(id)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.telemetry.Increment(ctx, "task_delete_total",
					attribute.String("status", "not_found"))
				s.telemetry.RecordDuration(ctx, "task_delete_duration_seconds",
					start,
					attribute.String("status", "not_found"))
				return errz.NewPrettyError(http.StatusNotFound, "not_found", "entity not found", err)
			}
			s.telemetry.Increment(ctx, "task_delete_total",
				attribute.String("status", "error"))
			s.telemetry.RecordDuration(ctx, "task_delete_duration_seconds",
				start,
				attribute.String("status", "error"))
			s.telemetry.RecordError(ctx, err)
			return errors.Wrap(err, "failed to check task ownership")
		}
		
		if userIdStr, ok := userId.(string); ok && existingTask.CreatedBy != userIdStr {
			s.telemetry.Increment(ctx, "task_delete_total",
				attribute.String("status", "forbidden"))
			s.telemetry.RecordDuration(ctx, "task_delete_duration_seconds",
				start,
				attribute.String("status", "forbidden"))
			return errz.NewPrettyError(http.StatusForbidden, "forbidden", "you don't have permission to delete this task", nil)
		}
	}

	_, err := s.q.WithContext(ctx).Task.Where(s.q.Task.ID.Eq(id)).Delete()
	if err != nil {
		s.telemetry.Increment(ctx, "task_delete_total",
			attribute.String("status", "error"))
		s.telemetry.RecordDuration(ctx, "task_delete_duration_seconds",
			start,
			attribute.String("status", "error"))
		s.telemetry.RecordError(ctx, err)
		return errors.Wrap(err, "failed to delete task")
	}

	// Record success
	s.telemetry.Increment(ctx, "task_delete_total",
		attribute.String("status", "success"))
	s.telemetry.RecordDuration(ctx, "task_delete_duration_seconds",
		start,
		attribute.String("status", "success"))
	return nil
}
