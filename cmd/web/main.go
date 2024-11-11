package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

type config struct {
	port int
	env  string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	cfg      config
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8000, "Server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	infologger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errorlogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := application{
		cfg:      cfg,
		infoLog:  infologger,
		errorLog: errorlogger,
	}

	r := chi.NewRouter()
	app.routes(r)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 8 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		app.errorLog.Println(err)
	}

}

