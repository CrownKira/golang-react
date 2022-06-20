package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// add routes() method to *application type; pointer type
// return router pointer type
// router: route url to func
// func (app *application) routes() *httprouter.Router {
// qn: why *httprouter.Router satisfies http.Handler?
// ans:
// ServeHTTP makes the router implement the http.Handler interface.
// func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {...} (check the source code)
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// this router than handle the paths below
	// method, path, handler
	// method: get method
	// listen to /status path
	// then pass to statusHandler to handle the request
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.getAllMoviesByGenre)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodPost, "/v1/admin/editmove", app.editmovie)

	return app.enableCORS(router)
}
