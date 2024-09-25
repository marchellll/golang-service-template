package service

import (
	"context"
	"golang-service-template/internal/dao/model"
	"golang-service-template/internal/errz"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type TaskService interface {
	Create(ctx context.Context, task model.Task) (*model.Task, error)
	Get(ctx context.Context, id string) (*model.Task, error)
	Find(ctx context.Context) ([]model.Task, error)
	Update(ctx context.Context, task model.Task) (*model.Task, error)
	Delete(ctx context.Context, id string) error
}

type taskService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewTaskService(i *do.Injector) (TaskService, error) {
	return &taskService{
		db:    do.MustInvoke[*gorm.DB](i),
		redis: do.MustInvoke[*redis.Client](i),
	}, nil
}

// Create implements TaskService.
func (s *taskService) Create(ctx context.Context, entity model.Task) (*model.Task, error) {
	entityp := &entity
	if err := s.db.WithContext(ctx).Create(entityp).Error; err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to create task", err)
	}
	return entityp, nil
}

// Get implements TaskService.
func (s *taskService) Get(ctx context.Context, id string) (*model.Task, error) {
	entity := &model.Task{}

	err := s.db.WithContext(ctx).First(entity, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errz.NewPrettyError(http.StatusNotFound, "not found", "entity not found", err)
	}

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entity", err)
	}

	return entity, nil
}

// GetAll implements TaskService.
func (s *taskService) Find(ctx context.Context) ([]model.Task, error) {
	var entities []model.Task
	if err := s.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entities", err)
	}

	return entities, nil
}

// Update implements TaskService.
func (s *taskService) Update(ctx context.Context, entity model.Task) (*model.Task, error) {
	err := s.db.WithContext(ctx).Model(&model.Task{}).Select("*").Updates(entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errz.NewPrettyError(http.StatusNotFound, "not found", "entity not found", err)
	}

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to update entities", err)
	}

	return &entity, nil
}

// Delete implements TaskService.
func (s *taskService) Delete(ctx context.Context, id string) error {
	if err := s.db.WithContext(ctx).Delete(&model.Task{}, id).Error; err != nil {
		return errors.Wrap(err, "failed to delete task")
	}
	return nil
}
