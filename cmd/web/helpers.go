package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
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

func (app *application) authenticatedUser(r *http.Request) int {
	session, err := app.session.Get(r, "user-session")
	if err != nil {
		app.errorLog.Println("Error retrieving session:", err)
		return 0
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return 0
	}

	return userID
}
