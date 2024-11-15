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
	err := ts.Execute(buffer, app.addDefaultData(td, r))
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	buffer.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}
