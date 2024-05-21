package controllers

import (
	"golang-service-template/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// interface
type HealthzController interface {
	Healthz() echo.HandlerFunc
}

// the struct that implements the interface
// and its dependencies
type healthzController struct {
	healthService services.HealthService
}

// New method
func NewHealthzController(healthService services.HealthService) HealthzController {
	return &healthzController{
		healthService: healthService,
	}
}

// Healthz method
func (controller *healthzController) Healthz() echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Message string `json:"message"`
		}

		body := new(req)

		if err := c.Bind(body); err != nil {
      return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

		requestBody := string(body.Message)
		response := controller.healthService.Healthcheck(requestBody)

		resp := req{
			Message: response,
		}

		c.JSON(http.StatusOK, resp)

		return nil

	}
}
