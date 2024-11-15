package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes(r *chi.Mux) {
	r.Use(app.recoverPanic)

	fs := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fs))
	r.Route("/", func(r chi.Router) {
		r.Get("/", app.homePage)
		r.Get("/login", app.loginPage)
		r.Get("/register", app.registerPage)
		r.Post("/register", app.registerUser)
	})
}
