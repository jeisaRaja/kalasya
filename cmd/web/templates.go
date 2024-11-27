package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/jeisaraja/kalasya/pkg/forms"
	"github.com/jeisaraja/kalasya/pkg/models"
)

func humanDate(t time.Time) string {
	return t.Format("04 Feb 2008")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type Error struct {
	StatusCode int
	Message    string
}

type templateData struct {
	AuthenticatedUser *models.User
	CurrentYear       int
	CSRFToken         string
	Form              *forms.Form
	Flash             string
	Error             Error
	Blog              *models.BlogView
	Post              *models.PostView
	Posts             []*models.Post
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
