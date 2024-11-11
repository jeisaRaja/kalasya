package main

import (
	"fmt"
	"net/http"
)

func (app *application) helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world!\n")
}
