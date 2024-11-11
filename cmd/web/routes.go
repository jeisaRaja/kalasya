package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes(r *chi.Mux) {
	fs := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fs))
	r.Route("/", func(r chi.Router) {
		r.Get("/login", app.loginPage)
	})
}
