package controllers

import (
	"golang-service-template/internal/services"
	"net/http"
)

// interface
type HealthzController interface {
	Healthz() http.Handler
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
func (c *healthzController) Healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type req struct {
			Message string `json:"message"`
		}

		var body req

		body, err := decode[req](r)
		// TODO: validate the request

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		requestBody := string(body.Message)
		response := c.healthService.Healthcheck(requestBody)

		resp := req{
			Message: response,
		}

		encode(w, r, http.StatusOK, resp)
	})
}
