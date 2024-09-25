package service

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/samber/do"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthService interface {
	Healthcheck(ctx context.Context) error
}

type healthService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthService(i *do.Injector) (HealthService, error) {
	return &healthService{
		db:    do.MustInvoke[*gorm.DB](i),
		redis: do.MustInvoke[*redis.Client](i),
	}, nil
}

func (service *healthService) Healthcheck(ctx context.Context) error {

	status := service.redis.Ping(ctx)
	if status.Err() != nil {
		return errors.Join(status.Err(), errors.New("failed to ping redis"))
	}

	result := service.redis.Get(ctx, "healthcheck")
	if result.Err() != nil && result.Err() != redis.Nil {
		return errors.Join(result.Err(), errors.New("failed to get healthcheck"))
	}

	if result.Val() == "OK" {
		// still OK, lets wait for expiry time before ping db again
		return nil
	}

	db, err := service.db.DB()
	if err != nil {
		return errors.Join(err, errors.New("failed to get db connection"))
	}

	err = db.Ping()
	if err != nil {
		return errors.Join(err, errors.New("failed to get ping db"))
	}

	service.redis.SetEx(ctx, "healthcheck", "OK", time.Second*30)

	return nil
}
