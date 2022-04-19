package main

import (
	"context"
	"errors"
	"net/http"
	"runtime"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hanna-yhchen/q-notes/internal/helpers"
	"github.com/hanna-yhchen/q-notes/internal/models"
	"github.com/justinas/nosurf"
)

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

// noteContext loads a note object based on the URL parameter and adds it to
// the context.
func noteContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helpers.NotFound(w)
			return
		}

		note, err := app.NoteModel.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				helpers.NotFound(w)
			} else {
				helpers.ServerError(w, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), helpers.ContextKeyNote, note)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// authenticate checks for auth status. If there is an active and authenticated
// user ID existed in the session, then records this status in the context passing
// through the follwing handler chain.
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if exists := app.Session.Exists(r, "authenticatedUserID"); !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.UserModel.Get(app.Session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) {
			app.Session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			helpers.ServerError(w, err)
			return
		}

		if !user.Active {
			app.Session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), helpers.ContextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireAuthentication redirects the client to the login page if it is not
// authenticated.
func requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Set the "Cache-Control: no-store" header so that any pages requiring
		// authentication are not stored in the client's cache.
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
