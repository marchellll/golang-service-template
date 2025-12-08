package app

import (
	"context"
	"golang-service-template/internal/common"
	"golang-service-template/internal/handler"
	"golang-service-template/internal/service"
	"io"

	"github.com/samber/do"

	"github.com/rs/zerolog"
)

// Container is a struct that holds all the components
// of the application: controllers, services, middlewares, auth, etc
type Container struct {
	// Logger
	Logger zerolog.Logger

	// configs
	Config common.Config

	// services
	HealthService service.HealthService
	service.TaskService

	// controllers
	HealthController handler.HealthzController
}

func NewInjector(
	ctx context.Context,
	getenv func(string) string,
	stdout, stderr io.Writer,
) *do.Injector {

	injector := do.New()

	// logger
	do.ProvideValue(injector, NewLogger(stdout))

	// configs
	do.ProvideValue(injector, NewConfig(getenv))

	// telemetry
	do.Provide(injector, NewTelemetry)

	do.Provide(injector, ConnectDB)
	do.Provide(injector, ConnectRedis)

	// services
	do.Provide(injector, service.NewHealthService)
	do.Provide(injector, service.NewTaskService)

	// handler
	do.Provide(injector, handler.NewHealthzController)
	do.Provide(injector, handler.NewTaskController)

	return injector
}
