package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string) {
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
	err := ts.Execute(buffer, nil)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	buffer.WriteTo(w)
}
