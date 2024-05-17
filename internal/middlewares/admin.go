package middlewares

import "net/http"

func AdminOnly(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement the logic

		// if !currentUser(r).IsAdmin {
		// 	http.NotFound(w, r)
		// 	return
		// }

		h.ServeHTTP(w, r)
	})
}