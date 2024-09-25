package handler

import (
	"golang-service-template/internal/service"
	"net/http"

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

// Healthz method
func (controller *healthzController) GetHealthz() echo.HandlerFunc {
	type req struct {
		Message string `json:"message"`
	}

	return func(c echo.Context) error {
		body := new(req)

		if err := c.Bind(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err := controller.healthService.Healthcheck(c.Request().Context())

		requestBody := string(body.Message)
		message := "I am healthty 🫡. This is your echo: " + requestBody + "."

		if err != nil {
			message = err.Error()
		}

		resp := req{
			Message: message,
		}

		c.JSON(http.StatusOK, resp)

		return nil

	}
}

// Healthz method
func (controller *healthzController) GetReadyz() echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Message string `json:"message"`
		}

		body := new(req)

		if err := c.Bind(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err := controller.healthService.Healthcheck(c.Request().Context())

		requestBody := string(body.Message)
		message := "I am healthty 🫡. This is your echo: " + requestBody + "."

		if err != nil {
			message = err.Error()
		}

		resp := req{
			Message: message,
		}

		c.JSON(http.StatusOK, resp)

		return nil

	}
}

func (controller *healthzController) Errorz() echo.HandlerFunc {
	return func(c echo.Context) error {
		return errors.New("errorz")
	}
}
