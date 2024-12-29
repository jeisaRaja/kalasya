package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gosimple/slug"
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

	var user models.UserRegistration
	err = form.GetInstance(&user)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	err = app.service.CreateUserWithBlog(&user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEmailDuplicate):
			form.Errors.Add("email", "email taken")
		case errors.Is(err, models.ErrSubdomainDuplicate):
			form.Errors.Add("subdomain", "subdomain taken")
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
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
	if err != nil {
		app.errorLog.Println("Failed to parse form data:", err)
		app.serverErrorResponse(w, r, err)
		return
	}
	form := forms.New(r.PostForm)
	id, err := app.service.AuthenticateUser(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		app.infoLog.Println("invalid credentials")
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.errorLog.Println("error when user model authenticate:", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	if id == nil {
		app.errorLog.Println("AuthenticateUser returned nil ID")
		app.serverErrorResponse(w, r, errors.New("server error: nil ID returned"))
		return
	}

	userSession, err := app.session.Get(r, "user-session")
	if err != nil {
		app.infoLog.Println("error when getting session")
		app.serverErrorResponse(w, r, err)
		return
	}

	app.infoLog.Println("id is ", *id)
	userSession.Values["user_id"] = *id
	err = userSession.Save(r, w)
	if err != nil {
		app.infoLog.Println("error when setting session value")
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) blogHomePage(w http.ResponseWriter, r *http.Request) {
	subdomain := chi.URLParam(r, "subdomain")
	post, err := app.service.GetBlogHome(subdomain)
	if err == nil {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if post.Content == "" {
		post.Content = "No Content Yet"
	}

	post.ContentHTML, err = app.toHTML(post.Content)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.render(w, r, "blogPost.page.tmpl", &templateData{
		Post: post,
	})
}

func (app *application) dashboardHomePage(w http.ResponseWriter, r *http.Request) {
	form := forms.New(url.Values{})
	user, ok := r.Context().Value(contextKeyUser).(*models.UserClient)
	if !ok {
		app.notFoundResponse(w, r)
		return
	}
	post, err := app.service.GetBlogHome(user.Subdomain)
	if err == models.ErrRecordNotFound {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	form.Add("homeContent", post.Content)
	app.render(w, r, "dashboardHome.page.tmpl", &templateData{Form: form})
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
	user, ok := r.Context().Value(contextKeyUser).(*models.UserClient)
	if !ok {
		app.notFoundResponse(w, r)
		return
	}
	post, err := app.service.GetBlogHome(user.Subdomain)
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
		return
	}
	form := forms.New(r.PostForm)
	value := form.Get("homeContent")
	post.Content = strings.TrimSpace(value)
	err = app.models.Post.Update(blog, post)
	err = app.service.UpdateBlogPost(posi)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/blog/%s/dashboard", blog.Subdomain), http.StatusSeeOther)
}

func (app *application) dashboardPostsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		app.notFoundResponse(w, r)
		return
	}
	posts, err := app.service.GetPosts(user.Subdomain)
	if err == models.ErrRecordNotFound {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.render(w, r, "dashboardPosts.page.tmpl", &templateData{Posts: posts})
}

func (app *application) dashboardCreatePostPage(w http.ResponseWriter, r *http.Request) {
	subdomain := chi.URLParam(r, "subdomain")
	form := forms.New(r.PostForm)
	app.render(w, r, "dashboardPost.page.tmpl", &templateData{
		Form: form,
		Post: &models.Post{},
		Blog: &models.Blog{Subdomain: subdomain},
	})
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		app.notFoundResponse(w, r)
		return
	}
	blog, _, err := app.models.Blogs.Get(user.Subdomain)
	app.infoLog.Printf("%#v", blog)
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
		return
	}
	form := forms.New(r.PostForm)

	publish := r.PostFormValue("publish") == "true"

	post := models.Post{
		BlogID:    blog.ID,
		Slug:      slug.Make(form.Get("title")),
		Title:     form.Get("title"),
		Content:   form.Get("content"),
		Published: publish,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = app.models.Post.CreatePost(&post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) dashboardPostPage(w http.ResponseWriter, r *http.Request) {
	postSlug := chi.URLParam(r, "post")
	if postSlug == "" {
		app.notFoundResponse(w, r)
		return
	}
	post, err := app.service.GetPost(postSlug, true)
	if err == models.ErrRecordNotFound {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	post.IsEdit = true

	form := forms.New(url.Values{})
	form.Add("title", post.Title)
	form.Add("content", post.Content)
	form.Add("published", strconv.FormatBool(post.Published))

	app.render(w, r, "dashboardPost.page.tmpl", &templateData{Form: form, Post: post})
}

func (app *application) blogPage(w http.ResponseWriter, r *http.Request) {
	subdomain := chi.URLParam(r, "subdomain")
	blog, _, err := app.models.Blogs.Get(subdomain)
	if err == models.ErrRecordNotFound {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	blog.NavHTML, err = app.parseBlogNav(blog.Nav, subdomain)

	app.render(w, r, "blogPosts.page.tmpl", &templateData{Blog: blog})
}
