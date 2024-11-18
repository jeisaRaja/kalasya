package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
		if app.authenticatedUser(r) == 0 {
			http.Redirect(w, r, "login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
