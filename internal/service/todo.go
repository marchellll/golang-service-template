package service

import (
	"context"
	"golang-service-template/internal/dao/model"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"gorm.io/gorm"
)



type TodoService interface {
	GetAll(ctx context.Context) ([]model.Todo, error)
	Create(ctx context.Context, todo model.Todo) (*model.Todo, error)
	Delete(ctx context.Context, id uint64) error
}

type todoService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewTodoService(i *do.Injector) (TodoService, error) {
	return &todoService{
		db: do.MustInvoke[*gorm.DB](i),
		redis: do.MustInvoke[*redis.Client](i),
	}, nil
}

// Create implements TodoService.
func (s *todoService) Create(ctx context.Context, todo model.Todo) (*model.Todo, error) {
	todop := &todo;
	if err := s.db.WithContext(ctx).Create(todop).Error; err != nil {
		return todop, errors.Wrap(err, "failed to create todo")
	}
	return todop, nil
}

// Delete implements TodoService.
func (s *todoService) Delete(ctx context.Context, id uint64) error {
	if err := s.db.WithContext(ctx).Delete(&model.Todo{}, id).Error; err != nil {
		return errors.Wrap(err, "failed to delete todo")
	}
	return nil
}


// GetAll implements TodoService.
func (s *todoService) GetAll(ctx context.Context) ([]model.Todo, error) {
	var todos []model.Todo
	if err := s.db.WithContext(ctx).Find(&todos).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get todos")
	}
	return todos, nil
}
