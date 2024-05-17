package internal

import "net/http"

func NewServer(di Container) http.Handler {
	mux := http.NewServeMux()

	regiterRoutes(mux, di)

	var handler http.Handler = mux
	// TODO: add global middleware here
	// handler = middleware1(handler)

	return handler
}
