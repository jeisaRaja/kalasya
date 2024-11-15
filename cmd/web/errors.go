package main

import (
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusInternalServerError, err)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err)
}
