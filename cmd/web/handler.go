package main

import (
	"errors"
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

	blogPost.ContentHTML, err = app.toHTML(blogPost.Content)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	blog.NavHTML, err = app.parseBlogNav(blog.Nav, blog.Subdomain)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.render(w, r, "blogPost.page.tmpl", &templateData{
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
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		app.notFoundResponse(w, r)
		return
	}
	blog, post, err := app.models.Blogs.Get(user.Subdomain)
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
	err = app.models.BlogPost.Update(blog, post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) dashboardPostsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		app.notFoundResponse(w, r)
		return
	}
	blogID, err := app.models.Blogs.GetID(user.Subdomain)
	if err == models.ErrRecordNotFound {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	posts, err := app.models.BlogPost.GetPosts(*blogID)
  app.infoLog.Println(posts)
  app.infoLog.Println(*blogID)
	app.render(w, r, "dashboardPosts.page.tmpl", &templateData{BlogPosts: posts})
}

func (app *application) dashboardCreatePostPage(w http.ResponseWriter, r *http.Request) {
	form := forms.New(r.PostForm)
	app.render(w, r, "dashboardPost.page.tmpl", &templateData{
		Form:     form,
		BlogPost: &models.BlogPost{},
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

	post := models.BlogPost{
		BlogID:    blog.ID,
		Slug:      slug.Make(form.Get("title")),
		Title:     form.Get("title"),
		Content:   form.Get("content"),
		Published: publish,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = app.models.BlogPost.Insert(&post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) dashboardEditPostPage(w http.ResponseWriter, r *http.Request) {
	postSlug := chi.URLParam(r, "post")
	if postSlug == "" {
		app.notFoundResponse(w, r)
		return
	}
	post, err := app.models.BlogPost.GetBySlug(postSlug)
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

	app.render(w, r, "dashboardPost.page.tmpl", &templateData{Form: form, BlogPost: post})
}
