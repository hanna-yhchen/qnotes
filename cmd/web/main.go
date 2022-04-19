package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/hanna-yhchen/q-notes/internal/config"
	"github.com/hanna-yhchen/q-notes/internal/handlers"
	"github.com/hanna-yhchen/q-notes/internal/helpers"
	"github.com/hanna-yhchen/q-notes/internal/render"
)

var (
	app    config.Application
	addr   = flag.String("addr", ":8080", "HTTP network address")
	secret = flag.String("secret", "Aof@fpaOEdAJepFls=(5&aBPeKOfjAk3", "Secret key for the session cookies")
	// TODO: Setup MySQL tables for demo mode.
	// demo := flag.Bool("demo", "false", "Demo mode")
)

func main() {
	flag.Parse()
	app = newApp()

	helpers.NewHelpers(&app)
	render.NewRenderer(&app)
	handlers.NewHandlers(&app)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  routes(),
	}

	app.InfoLog.Printf("Starting server on %s", *addr)
	app.ErrorLog.Fatal(srv.ListenAndServe())
}

func newApp() config.Application {
	errorLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	// session.Secure = true

	app := config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Session:  session,
	}

	tc, err := render.NewTemplateCache("./ui/template/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app.TemplateCache = tc

	return app
}
