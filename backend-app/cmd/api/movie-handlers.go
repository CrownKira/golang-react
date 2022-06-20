package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	// get the params from context
	// provide the context: context is found in request
	params := httprouter.ParamsFromContext(r.Context())

	// convert str to int
	// extract id param
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	// app.logger.Println("id is", id)

	// get from the app.models
	// facade thinking
	fmt.Println("getting movie", id)
	movie, err := app.models.DB.Get(id)

	// movie := models.Movie{
	// 	ID:          id,
	// 	Title:       "Some movie",
	// 	Description: "Some description",
	// 	Year:        2021,
	// 	ReleaseDate: time.Date(2021, 01, 01, 01, 0, 0, 0, time.Local),
	// 	Runtime:     100,
	// 	Rating:      5,
	// 	MPAARating:  "PG-13",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }

	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// write to response
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GenresAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) getAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get genre id from param
	genreID, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// get the movie based on the genre id
	movies, err := app.models.DB.All(genreID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// write movie json to response
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

// all modifying operations
func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {

}

func (app *application) insertMovie(w http.ResponseWriter, r *http.Request) {

}

func (app *application) editmovie(w http.ResponseWriter, r *http.Request) {
	type jsonRep struct {
		OK boolean `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}

	err := app.writeJSON(w, http.s)
}

func (app *application) searchMovie(w http.ResponseWriter, r *http.Request) {

}
