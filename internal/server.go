package internal

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func NewServer(di Container) *echo.Echo {
	e := echo.New()

  e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	// healthz
	e.GET("/healthz", di.HealthController.Healthz())
	e.POST("/healthz", di.HealthController.Healthz())

	return e
}
