package service

import (
	"context"
	"golang-service-template/internal/dao/model"
	"golang-service-template/internal/dao/query"
	"golang-service-template/internal/errz"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
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
	db    *gorm.DB
	q     *query.Query
	redis *redis.Client
}

func NewTaskService(i *do.Injector) (TaskService, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &taskService{
		db:    db,
		q:     query.Use(db),
		redis: do.MustInvoke[*redis.Client](i),
	}, nil
}

// Create implements TaskService.
func (s *taskService) Create(ctx context.Context, entity model.Task) (*model.Task, error) {
	newID, err := uuid.NewV7()

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to generate new id", err)
	}

	entityp := &entity
	entityp.ID = newID.String()
	entityp.CreatedBy = newID.String() // TODO: get user id from context
	if err := query.Use(s.db).WithContext(ctx).Task.Create(entityp); err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to create task", err)
	}
	return entityp, nil
}

// Get implements TaskService.
func (s *taskService) Get(ctx context.Context, id string) (*model.Task, error) {
	entity, err := s.q.WithContext(ctx).Task.Where(s.q.Task.ID.Eq(id)).First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errz.NewPrettyError(http.StatusNotFound, "not_found", "entity not found", err)
	}

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entity", err)
	}

	return entity, nil
}

// GetAll implements TaskService.
func (s *taskService) Find(ctx context.Context) ([]*model.Task, error) {
	entities, err := s.q.WithContext(ctx).Task.Find()

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entities", err)
	}

	return entities, nil
}

// GetAll implements TaskService.
func (s *taskService) FindByUserId(ctx context.Context, userId string) ([]*model.Task, error) {
	entities, err := s.q.WithContext(ctx).Task.Where(s.q.Task.CreatedBy.Eq(userId)).Find()

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entities", err)
	}

	return entities, nil
}

// Update implements TaskService.
// using map here to avoid headache of handling Go's zero value
// we pass whatever passed validation in handler
func (s *taskService) Update(ctx context.Context, id string, entity map[string]any) (*model.Task, error) {
	_, err := s.q.WithContext(ctx).Task.Where(s.q.Task.ID.Eq(id)).Updates(entity)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errz.NewPrettyError(http.StatusNotFound, "not_found", "entity not found", err)
	}

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to update entities", err)
	}

	return s.Get(ctx, id)
}

// Delete implements TaskService.
func (s *taskService) Delete(ctx context.Context, id string) error {
	_, err := s.q.WithContext(ctx).Task.Where(s.q.Task.ID.Eq(id)).Delete()
	if err != nil {
		return errors.Wrap(err, "failed to delete task")
	}
	return nil
}
