package app

import (
	"github.com/samber/do"

	"golang-service-template/internal/common"
	"golang-service-template/internal/errz"
	"golang-service-template/internal/handler"
	"golang-service-template/internal/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// routes and middlewares here
func addRoutes(
	e *echo.Echo,
	injector *do.Injector,
) {
	logger := do.MustInvoke[zerolog.Logger](injector)

	// global middlewares
	e.Use(echo_middleware.Recover())
	e.Use(errz.ErrzMiddleware())
	e.Use(echo_middleware.CORS())
	e.Use(middleware.LoggerMiddleware(logger))

	// routes
	addHealthzRoutes(injector, e)
	addTaskRoutes(injector, e)

	// root route
	e.Any("/", echo.WrapHandler(http.NotFoundHandler()))
}

func addHealthzRoutes(injector *do.Injector, e *echo.Echo) {
	healthController := do.MustInvoke[handler.HealthzController](injector)
	e.GET("/healthz", healthController.GetHealthz())
	e.POST("/healthz", healthController.GetHealthz())
	e.GET("/readyz", healthController.GetReadyz())
	e.GET("/errorz", healthController.Errorz())
}

func addTaskRoutes(injector *do.Injector, e *echo.Echo) {
	taskGroup := e.Group("/tasks")

	// taskGroup.Use(echo_middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == "username" && password == "password" {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	taskGroup.POST("", do.MustInvoke[handler.TaskController](injector).Create())
	taskGroup.GET("", do.MustInvoke[handler.TaskController](injector).Find())
	taskGroup.GET("/:id", do.MustInvoke[handler.TaskController](injector).GetById())
	taskGroup.PATCH("/:id", do.MustInvoke[handler.TaskController](injector).Update())
	taskGroup.DELETE("/:id", do.MustInvoke[handler.TaskController](injector).Delete())

	securedTaskGroup := e.Group("/secured/tasks")
	securedTaskGroup.Use(middleware.ValidateJWTMiddleware(
		do.MustInvoke[common.Config](injector),
		do.MustInvoke[zerolog.Logger](injector),
	))

	securedTaskGroup.GET("", do.MustInvoke[handler.TaskController](injector).FindByUserId())

}
