package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func (app *application) routes(r *chi.Mux) {
	csrfMiddleware := csrf.Protect(
		[]byte(os.Getenv("CSRF_KEY")),
		csrf.HttpOnly(true),
    csrf.Secure(false),
    csrf.Path("/"),
    csrf.FieldName("csrf_token"),
	)

	r.Use(app.recoverPanic)
	fs := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fs))
	r.Route("/", func(r chi.Router) {
		r.Use(csrfMiddleware)
		r.Get("/", app.homePage)
		r.Get("/login", app.loginPage)
		r.Post("/login", app.loginUser)
		r.Post("/logout", app.logoutUser)
		r.Get("/register", app.registerPage)
		r.Post("/register", app.registerUser)
		r.With(app.requireAuthenticatedUser).Get("/dashboard", app.dashboardPage)
		r.Get("/blog/{subdomain}", app.subdomainHandler)
	})
}
