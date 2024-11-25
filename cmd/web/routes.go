package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		r.Use(middleware.StripSlashes)
		r.Use(csrfMiddleware)
		r.Use(app.authenticate)

		// Check authentication only when necessary, not globally
		r.With(app.redirectIfAuthenticated).Get("/login", app.loginPage)
		r.With(app.redirectIfAuthenticated).Get("/register", app.registerPage)

		// Routes for login, logout and register
		r.Post("/login", app.loginUser)
		r.Post("/logout", app.logoutUser)
		r.Post("/register", app.registerUser)

		// Public routes
		r.Get("/", app.homePage)
		r.Route("/blog/{subdomain}", func(r chi.Router) {
			r.Get("/", app.blogHomePage)
			r.Get("/{nav}", app.blogPage)
		})

		// Authenticated routes
		r.Route("/blog/{subdomain}/dashboard", func(r chi.Router) {
			r.Use(app.requireAuthenticatedUser)
			r.Use(app.requireAuthorizedUser)

			r.Get("/", app.dashboardPage)
			r.Get("/posts", app.dashboardPostsPage)
			r.Get("/posts/{post}", app.dashboardPostPage)
			r.Post("/posts", app.createPost)
			r.Get("/create-post", app.dashboardCreatePostPage)
			r.Post("/home", app.updateBlogHome)
		})
	})

	// This route will handle unmatched paths
	r.NotFound(app.notFoundResponse)
}
