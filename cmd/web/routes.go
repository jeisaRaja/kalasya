package main

import "github.com/go-chi/chi/v5"

func (app *application) routes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/login", app.loginPage)
	})
}
