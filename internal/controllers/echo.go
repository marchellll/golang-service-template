package controllers

import (
	"golang-service-template/internal/services"
	"net/http"
)

// interface
type EchoController interface {
	Echo() (http.Handler)
}


// the struct that implements the interface
// and its dependencies
type echoController struct {
	echoService services.EchoService
}

// New method
func NewEchoController(echoService services.EchoService) EchoController {
	return &echoController{
		echoService: echoService,
	}
}


// Echo method
func (ec *echoController) Echo() (http.Handler) {
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

		toEcho := string(body.Message)
		toEcho = ec.echoService.Echo(toEcho)

		resp := req{
			Message: toEcho,
		}

		encode(w, r, http.StatusOK, resp)
	})
}