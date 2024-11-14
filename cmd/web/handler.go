package main

import (
	"net/http"
)

func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl")
}

func (app *application) homePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.tmpl")
}

