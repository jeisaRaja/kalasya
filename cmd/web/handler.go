package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jeisaraja/kalasya/pkg/forms"
	"github.com/jeisaraja/kalasya/pkg/models"
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
	fmt.Println("home page")
	app.render(w, r, "home.page.tmpl", nil)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	form := forms.New(r.PostForm)
	models.ValidateUser(form)

	if !form.Valid() {
		app.render(w, r, "register.page.tmpl", &templateData{Form: form})
		return
	}

	var user models.User
	err = form.GetInstance(&user)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	err = app.models.Users.Exists(&user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEmailDuplicate):
			form.Errors.Add("email", "email taken")
		case errors.Is(err, models.ErrSubdomainDuplicate):
			form.Errors.Add("subdomain", "subdomain taken")
		}
		app.render(w, r, "register.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		app.errorLog.Println(err)
		app.render(w, r, "register.page.tmpl", &templateData{Form: form})
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Creating a new user...\nUser: %+v", user)))
}

func (app *application) subdomainHandler(w http.ResponseWriter, r *http.Request) {
	subdomain := chi.URLParam(r, "subdomain")
	blog, err := app.models.Blogs.Get(subdomain)
	if err != nil {
    app.errorLog.Println(err)
		app.notFoundResponse(w, r)
		return
	}

  fmt.Println("blog is: ", blog)
	fmt.Fprintf(w, "%#v", blog)
}
