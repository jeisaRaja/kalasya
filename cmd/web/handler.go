package main

import (
	"net/http"

	"github.com/jeisaraja/kalasya/pkg/forms"
)

func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", nil)
}

func (app *application) registerPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) homePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.tmpl", nil)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("blogtitle", "subdomain", "email", "password")
	form.MaxLength("blogtitle", 100)
	form.MaxLength("subdomain", 50)
	form.MinLength("blogtitle", 10)
	form.MinLength("subdomain", 3)
	form.MinLength("password", 8)
	form.MaxLength("password", 30)
	form.EmailValid("email")

	if !form.Valid() {
		app.render(w, r, "register.page.tmpl", &templateData{Form: form})
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Creating a new user..."))
}
