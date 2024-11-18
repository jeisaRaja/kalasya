package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/jeisaraja/kalasya/pkg/models"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	cfg           config
	templateCache map[string]*template.Template
	models        models.Models
	session       *sessions.CookieStore
}

func main() {
	var cfg config

	defaultDSN := "postgres://name:password@localhost:5433/kalasya?sslmode=disable"

	if dsnEnv := os.Getenv("DSN"); dsnEnv != "" {
		defaultDSN = dsnEnv
	}

	flag.IntVar(&cfg.port, "port", 8000, "Server port")
	flag.StringVar(&cfg.db.dsn, "dsn", defaultDSN, "Data source name")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	infologger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errorlogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorlogger.Fatalf("failed to create new template cache: %v", err)
	}

	db, err := openDB(cfg)
	if err != nil {
		errorlogger.Fatal(err)
	}

	infologger.Printf("connected to database at %s", cfg.db.dsn)

	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("AUTH_KEY")), []byte(os.Getenv("ENCRYPT_KEY")))

	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	app := application{
		cfg:           cfg,
		infoLog:       infologger,
		errorLog:      errorlogger,
		templateCache: templateCache,
		models:        models.New(db),
		session:       sessionStore,
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

	infologger.Printf("running on port :%d\n", cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		app.errorLog.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
