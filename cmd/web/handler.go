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
	session, err := app.session.Get(r, "user-session")
	if err != nil {
		app.errorLog.Println("failed to get user-session", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	flash, ok := session.Values["flash"].(string)
	if !ok {
		flash = ""
	}
	app.infoLog.Println(flash)

	delete(session.Values, "flash")
	err = session.Save(r, w)
	if err != nil {
		app.errorLog.Println("failed to save user-session", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.infoLog.Println("before rendering login.page.tmpl")
	app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil), Flash: flash})
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
	models.ValidateUserRegistration(form)

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

	session, err := app.session.Get(r, "user-session")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	session.Values["flash"] = "Registration successfull, please login to access your account."
	err = session.Save(r, w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) subdomainHandler(w http.ResponseWriter, r *http.Request) {
	subdomain := chi.URLParam(r, "subdomain")
	blog, err := app.models.Blogs.Get(subdomain)
	if err != nil {
		app.errorLog.Println(err)
		app.notFoundResponse(w, r)
		return
	}

	fmt.Fprintf(w, "%#v", blog)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	form := forms.New(r.PostForm)
	id, err := app.models.Users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.errorLog.Println("error when user model authenticate:", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	userSession, err := app.session.Get(r, "user-session")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	userSession.Values["user_id"] = id

	app.infoLog.Println("redirect to home")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
