package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	// handler: request -> response
	// handler has serveHttp method

	// invoke serveHttp method
	// set the header of the response
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add allow origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
