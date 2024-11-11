package main

import (
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
