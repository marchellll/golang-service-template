package internal

import (
	"golang-service-template/internal/controllers"
	"golang-service-template/internal/services"
)

// Container is a struct that holds all the components
// of the application: controllers, services, middlewares, auth, etc
type Container struct {
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
	// configs
	config := NewConfig()

	// services
	healthService := services.NewHealthService()

	// controllers
	healthController := controllers.NewHealthzController(healthService)

	return Container{
		Config: config,

		// services
		HealthService: healthService,

		// controllers
		HealthController: healthController,
	}
}
