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
	EchoService services.EchoService

	// controllers
	EchoController controllers.EchoController
}

// NewDI is a constructor for DI
// this will be the default DI for the application
// Test can create a new DI with their own (mini) dependencies
func NewContainer() Container {
	// configs
	config := NewConfig()

	// services
	echoService := services.NewEchoService()

	// controllers
	echoController := controllers.NewEchoController(echoService)

	return Container{
		Config: config,

		// services
		EchoService: echoService,

		// controllers
		EchoController: echoController,
	}
}
