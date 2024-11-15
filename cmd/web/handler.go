package main

import (
	"fmt"
	"net/http"

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

	var user models.User
	var blog models.Blog
	err = form.GetInstance(&user)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	err = form.GetInstance(&blog)
	if err != nil {
		app.errorLog.Println(err)
	}

	blogExists, err := app.models.Blogs.Exists(blog.Subdomain)
	if err != nil {
		app.errorLog.Println(err)
		app.render(w, r, "register.page.tmpl", &templateData{Form: form})
		return
	}

	if blogExists {
		form.Errors.Add("subdomain", "blog with this subdomain exists already")
		app.render(w, r, "register.page.tmpl", &templateData{Form: form})
		return
	}

	emailExists, err := app.models.Users.Exists(user.Email)
	if err != nil {
		app.errorLog.Println(err)
		app.render(w, r, "register.page.tmpl", &templateData{Form: form})
		return
	}

	if emailExists {
    app.infoLog.Println("email duplicate")
		form.Errors.Add("email", "user with this email exists already")
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
