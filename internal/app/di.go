package app

import (
	"golang-service-template/internal/controllers"
	"golang-service-template/internal/services"

	"go.uber.org/zap"
)

// Container is a struct that holds all the components
// of the application: controllers, services, middlewares, auth, etc
type Container struct {
	// Logger
	Logger *zap.Logger

	// configs
	Config Config

	// services
	HealthService services.HealthService

	// controllers
	HealthController controllers.HealthzController
}

// NewDI is a constructor for DI
// this will be the default DI for the application
// Test can create a new DI with their own (mini) dependencies
func NewContainer() Container {
	// logger
	logger := NewLogger()

	// configs
	config := NewConfig()

	gormDB := ConnectDB(logger, config)
	redis := ConnectRedis(logger, config)

	// services
	healthService := services.NewHealthService(gormDB, redis)

	// controllers
	healthController := controllers.NewHealthzController(healthService)

	return Container{
		// logger
		Logger: logger,

		Config: config,

		// services
		HealthService: healthService,

		// controllers
		HealthController: healthController,
	}
}
