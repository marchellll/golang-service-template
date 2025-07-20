package handler

import (
	"golang-service-template/internal/service"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/samber/do"

	"github.com/labstack/echo/v4"
)

// interface
type HealthzController interface {
	GetHealthz() echo.HandlerFunc
	GetReadyz() echo.HandlerFunc
	Errorz() echo.HandlerFunc
}

// the struct that implements the interface
// and its dependencies
type healthzController struct {
	healthService service.HealthService
}

// New method
func NewHealthzController(i *do.Injector) (HealthzController, error) {
	return &healthzController{
		healthService: do.MustInvoke[service.HealthService](i),
	}, nil
}

// GetHealthz - Liveness probe endpoint
// This should be fast and lightweight, only checking if the app is alive
func (controller *healthzController) GetHealthz() echo.HandlerFunc {
	type response struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}

	return func(c echo.Context) error {
		err := controller.healthService.LivenessCheck(c.Request().Context())

		if err != nil {
			resp := response{
				Status:    "unhealthy",
				Message:   "Liveness check failed: " + err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			}
			return c.JSON(http.StatusServiceUnavailable, resp)
		}

		resp := response{
			Status:    "healthy",
			Message:   "I am alive 🫡",
			Timestamp: time.Now().Format(time.RFC3339),
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// GetReadyz - Readiness probe endpoint
// This can be more thorough, checking if the app is ready to serve traffic
func (controller *healthzController) GetReadyz() echo.HandlerFunc {
	type response struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}

	return func(c echo.Context) error {
		err := controller.healthService.ReadinessCheck(c.Request().Context())

		if err != nil {
			resp := response{
				Status:    "not_ready",
				Message:   "Readiness check failed: " + err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			}
			return c.JSON(http.StatusServiceUnavailable, resp)
		}

		resp := response{
			Status:    "ready",
			Message:   "I am ready to serve traffic 🚀",
			Timestamp: time.Now().Format(time.RFC3339),
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func (controller *healthzController) Errorz() echo.HandlerFunc {
	return func(c echo.Context) error {
		return errors.New("errorz")
	}
}
