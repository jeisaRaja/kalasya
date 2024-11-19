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
    r.Use(app.authenticate)
		
		// Check authentication only when necessary, not globally
		r.With(app.redirectIfAuthenticated).Get("/login", app.loginPage)
		r.With(app.redirectIfAuthenticated).Get("/register", app.registerPage)

		// Routes for login, logout and register
		r.Post("/login", app.loginUser)
		r.Post("/logout", app.logoutUser)
		r.Post("/register", app.registerUser)

		// Authenticated routes
		r.With(app.requireAuthenticatedUser).Get("/dashboard", app.dashboardPage)

		// Public routes
		r.Get("/", app.homePage)
		r.Get("/blog/{subdomain}", app.blogHomePage)
	})

	// This route will handle unmatched paths
	r.NotFound(app.notFoundResponse)
}
