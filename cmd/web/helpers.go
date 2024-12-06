package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jeisaraja/kalasya/pkg/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	buffer := new(bytes.Buffer)
	ts, ok := app.templateCache[name]
	if !ok {
		app.errorLog.Printf("failed to get template, name: %s, templateCache: %+v", name, app.templateCache)

		for name, tmpl := range app.templateCache {
			app.errorLog.Printf("Template name: %s, Template contents: %+v\n", name, *tmpl)
		}
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("Template %s doesn't exist!", name))
		return
	}
	err := ts.Execute(buffer, app.addDefaultData(td, r, w))
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	buffer.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request, w http.ResponseWriter) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()

	td.CSRFToken = csrf.Token(r)

	blogInfo := r.Context().Value(contextKeyBlog).(*models.BlogView)
	td.Blog = blogInfo

	session, err := app.session.Get(r, "user-session")
	if err != nil {
		td.Flash = ""
	} else {
		flashes := session.Flashes()
		if len(flashes) > 0 {
			td.Flash, _ = flashes[0].(string)
		}
		session.Save(r, w)
	}
	return td
}

func (app *application) authenticatedUser(r *http.Request) *models.UserClient {
	user, ok := r.Context().Value(contextKeyUser).(*models.UserClient)
	if !ok {
		return nil
	}

	return user
}

func (app *application) authorizedUser(r *http.Request) bool {
	user, ok := r.Context().Value(contextKeyUser).(*models.UserClient)
	if !ok {
		return false
	}
	urlSubdomain := chi.URLParam(r, "subdomain")
	app.infoLog.Println(urlSubdomain)
	if user.Subdomain != urlSubdomain {
		return false
	}
	return true
}

func (app *application) toHTML(md string) (template.HTML, error) {
	var unsafe bytes.Buffer
	if err := goldmark.Convert([]byte(md), &unsafe); err != nil {
		return "", err
	}
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe.Bytes())

	return template.HTML(html), nil
}

func (app *application) parseBlogNav(nav, subdomain string) (template.HTML, error) {
	var unsafe strings.Builder

	if err := goldmark.Convert([]byte(nav), &unsafe); err != nil {
		return "", err
	}
	baseUrl := fmt.Sprintf("/blog/%s", subdomain)

	htmlStr := unsafe.String()
	htmlStr = strings.ReplaceAll(htmlStr, `href="/`, fmt.Sprintf(`href="%s/`, baseUrl))

	html := bluemonday.UGCPolicy().Sanitize(htmlStr)

	return template.HTML(html), nil
}
