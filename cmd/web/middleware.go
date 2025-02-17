package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jeisaraja/kalasya/pkg/models"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s\n", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) subdomainMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		subdomain := strings.Split(host, ".")[0]

		ctx := context.WithValue(r.Context(), "subdomain", subdomain)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) addBlogInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subdomain := chi.URLParam(r, "subdomain")
		if subdomain == "" {
			app.notFoundResponse(w, r)
			return
		}
		blogInfo, err := app.service.GetBlogInfo(subdomain)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyBlog, blogInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.session.Get(r, "user-session")
		userID, exists := session.Values["user_id"].(int)
		if !exists {
			app.infoLog.Println("session not exist")
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.service.GetAuthenticatedUser(userID)
		if err == models.ErrRecordNotFound {
			delete(session.Values, "user_id")
			session.Save(r, w)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) redirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthorizedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorized := app.authorizedUser(r)
		if !authorized {
			app.notFoundResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
