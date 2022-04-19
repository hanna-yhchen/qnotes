package main

import (
	"context"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hanna-yhchen/q-notes/internal/helpers"
	"github.com/hanna-yhchen/q-notes/internal/models"
	"github.com/justinas/nosurf"
)

type contextKey string

const contextKeyNote = contextKey("note")

// logger returns a logger handler using a custom LogFormatter.
func logger() func(http.Handler) http.Handler {
	color := true
	if runtime.GOOS == "windows" {
		color = false
	}
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: app.InfoLog, NoColor: !color})
}

// noSurf adds csrf protection to the routing chain.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
	})
	return csrfHandler
}

// noteContext loads a note object from the URL parameter and adds it to the context.
func noteContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var note *models.Note

		if id := chi.URLParam(r, "id"); id != "" {
			// note, err := app.noteModel.Get(id)
			// if err != nil {
			// 	if errors.Is(err, models.ErrNoRecord) {
			// 		app.notFound(w)
			// 	} else {
			// 		app.serverError(w, err)
			// 	}
			// 	return
			// }
		} else {
			helpers.NotFound(w)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyNote, note)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
