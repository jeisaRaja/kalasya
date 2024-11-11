package main

import (
	"fmt"
	"net/http"
)

func (app *application) loginPage(w http.ResponseWriter, r *http.Request) {
  fmt.Println("before calling render")
	app.render(w, r, "login", &templateData{})
}
