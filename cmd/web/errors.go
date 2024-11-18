package main

import (
	"errors"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message error) {
	w.WriteHeader(status)
	app.render(w, r, "error.page.tmpl", &templateData{Error: Error{StatusCode: status, Message: message.Error()}})
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusInternalServerError, err)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, errors.New("Not Found"))
}
