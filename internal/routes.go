package internal

import (
	"golang-service-template/internal/middlewares"
	"net/http"
)

func regiterRoutes(mux *http.ServeMux, di Container) {
	// handle all routes here

	mux.Handle("/healthz", middlewares.AdminOnly(di.HealthController.Healthz()))
}
