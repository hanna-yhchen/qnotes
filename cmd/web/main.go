package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/hanna-yhchen/q-notes/internal/config"
	"github.com/hanna-yhchen/q-notes/internal/handlers"
	"github.com/hanna-yhchen/q-notes/internal/helpers"
	"github.com/hanna-yhchen/q-notes/internal/models/mysql"
	"github.com/hanna-yhchen/q-notes/internal/render"
)

var (
	app    *config.Application
	addr   = flag.String("addr", ":8080", "HTTP network address")
	secret = flag.String("secret", "Aof@fpaOEdAJepFls=(5&aBPeKOfjAk3", "Secret key for the session cookies")
	dsn    = flag.String("dsn", "web:pass@/qnotes?parseTime=true", "MySQL data source name")
)

func main() {
	srv, closeDB := setupApp()
	defer closeDB()

	app.InfoLog.Printf("Starting server on %s", *addr)
	app.ErrorLog.Fatal(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
}

// setupApp setups app-wide config. Return a http.server and a function to close
// database connection.
func setupApp() (*http.Server, func()) {
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	tc, err := render.NewTemplateCache("./ui/template/")
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	app = &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Session:  session,
		TemplateCache: tc,
		NoteModel: &mysql.NoteModel{DB: db},
		UserModel: &mysql.UserModel{DB: db},
	}

	helpers.NewHelpers(app)
	render.NewRenderer(app)
	handlers.NewHandlers(app)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  routes(),
		TLSConfig: tlsConfig,
	}

	return srv, func() {
		db.Close()
	}
}

// openDB wraps sql.Open and returns a sql.DB connection pool for the given DSN.
func openDB(dsn string) (*sql.DB, error) {
	if db, err := sql.Open("mysql", dsn); err != nil {
		return nil, err
		// Ping creates and verifies a connection.
	} else if err = db.Ping(); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}