package main

import (
	"fmt"
	"net/http"
	"time"

	"movie.alvintanoto.id/internal/data"
	"movie.alvintanoto.id/internal/validator"
)

// Add createMovieHandler
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Create a new movie")
	// declare an anonymous struct
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// initialize new decoder
	// then use Decode() method to decode the body contents into the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// initialize validator instance
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	fmt.Fprintf(w, "%+v\n", input)
}

// Add showMovieHandler
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		// http.NotFound(w, r)

		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Batman",
		Runtime:   102,
		Genres:    []string{"vengeance", "suspense", "superhero"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// app.logger.Println(err)
		// http.Error(w, "The server cannot process your request", http.StatusInternalServerError)

		app.serverErrorResponse(w, r, err)
	}
}
