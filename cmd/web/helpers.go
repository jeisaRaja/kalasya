package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	buffer := bytes.Buffer{}
	ts, ok := app.templateCache[name]
	if !ok {
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("Template %s doesn't exist!", name))
	}
  fmt.Println("before ts.execute")
	err := ts.Execute(&buffer, td)
  fmt.Println("after ts.execute")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	buffer.WriteTo(w)
}
