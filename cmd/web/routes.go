package main

import "github.com/go-chi/chi/v5"

func (app *application) routes(r *chi.Mux) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/", app.helloworld)
	})
}
