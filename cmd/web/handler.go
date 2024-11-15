package main

import (
	"fmt"
	"net/http"
)

func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl")
}

func (app *application) registerPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.tmpl")
}

func (app *application) homePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.tmpl")
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	blogTitle := r.PostForm.Get("blog-title")
	email := r.PostForm.Get("email")
	subdomain := r.PostForm.Get("subdomain")
	password := r.PostForm.Get("password")

	response := fmt.Sprintf("Blog Title: %s\nEmail: %s\nSubdomain: %s\nPassword: %s", blogTitle, email, subdomain, password)

  fmt.Println(response)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
