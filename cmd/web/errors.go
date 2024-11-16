package main

import (
	"fmt"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %s", message)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusInternalServerError, err)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, nil)
}
