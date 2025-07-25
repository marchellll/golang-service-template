package service

import (
	"context"
	"fmt"
	"golang-service-template/internal/common"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/samber/do"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthService interface {
	LivenessCheck(ctx context.Context) error
	ReadinessCheck(ctx context.Context) error
}

type healthService struct {
	db         *gorm.DB
	redis      *redis.Client
	config     common.Config
	instanceID string
}

func NewHealthService(i *do.Injector) (HealthService, error) {
	instanceID := uuid.New().String()

	return &healthService{
		db:         do.MustInvoke[*gorm.DB](i),
		redis:      do.MustInvoke[*redis.Client](i),
		config:     do.MustInvoke[common.Config](i),
		instanceID: instanceID,
	}, nil
}

func (service *healthService) ReadinessCheck(ctx context.Context) error {

	status := service.redis.Ping(ctx)
	if status.Err() != nil {
		return errors.Join(status.Err(), errors.New("failed to ping redis"))
	}

	healthcheckKey := fmt.Sprintf("%s:healthcheck:%s", service.config.ServiceName, service.instanceID)
	result := service.redis.Get(ctx, healthcheckKey)
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

	timeout := time.Duration(service.config.HealthcheckTimeoutSeconds) * time.Second
	service.redis.SetEx(ctx, healthcheckKey, "OK", timeout)

	return nil
}

// LivenessCheck performs a basic liveness check - just returns current time (very lightweight)
// This indicates the application is alive and responding
func (service *healthService) LivenessCheck(ctx context.Context) error {
	// Simply return nil - the fact that we can execute this function means the app is alive
	// The handler will include the current timestamp in the response
	return nil
}
