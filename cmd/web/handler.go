package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jeisaraja/kalasya/pkg/forms"
	"github.com/jeisaraja/kalasya/pkg/models"
)

func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil)})
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

	session.AddFlash("Registration successfull, please login to access your account.")
	err = session.Save(r, w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	err = userSession.Save(r, w)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) blogHomePage(w http.ResponseWriter, r *http.Request) {
	subdomain := chi.URLParam(r, "subdomain")
	blog, blogPost, err := app.models.Blogs.Get(subdomain)
	if err != nil {
		app.errorLog.Println(err)
		app.notFoundResponse(w, r)
		return
	}
	if blogPost.Content == "" {
		blogPost.Content = "No Content Yet"
	}

	app.render(w, r, "post.page.tmpl", &templateData{
		Blog:     blog,
		BlogPost: blogPost,
	})
}

func (app *application) dashboardPage(w http.ResponseWriter, r *http.Request) {
	form := forms.New(url.Values{})
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		app.notFoundResponse(w, r)
		return
	}
	_, post, err := app.models.Blogs.Get(user.Subdomain)
	if err == models.ErrRecordNotFound {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.errorLog.Println("content", post.Content)
	form.Add("homeContent", post.Content)
	app.render(w, r, "dashboard.page.tmpl", &templateData{Form: form})
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	session, err := app.session.Get(r, "user-session")
	if err != nil {
		app.errorLog.Println("cannot get user-session:", err)
		app.serverErrorResponse(w, r, err)
	}

	delete(session.Values, "user_id")
	session.AddFlash("You've been logged out successfully.")
	err = session.Save(r, w)
	if err != nil {
		app.errorLog.Println("cannot save session:", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) updateBlogHome(w http.ResponseWriter, r *http.Request) {
	subdomain := chi.URLParam(r, "subdomain")

	blog, post, err := app.models.Blogs.Get(subdomain)
	if err == models.ErrRecordNotFound {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	form := forms.New(r.PostForm)
	value := form.Get("homeContent")
	post.Content = strings.TrimSpace(value)
	err = app.models.BlogPost.Update(blog, post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	app.render(w, r, "dashboard.page.tmpl", &templateData{Form: form})
}
