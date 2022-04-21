package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hanna-yhchen/q-notes/internal/handlers"
)

func routes() http.Handler {
	r := chi.NewRouter()

	baseChain := chi.Chain(
		logger(),
		middleware.RequestID,
		middleware.Recoverer,
		// prevent rendering if an XSS attack is detected
		middleware.SetHeader("X-XSS-Protection", "1; mode=block"),
		// disallow embedded into other sites to avoid click-jacking attacks
		middleware.SetHeader("X-Frame-Options", "deny"),
	)

	r.Use(app.Session.Enable, noSurf, authenticate)

	r.Get("/", handlers.Home)
	r.With(requireAuthentication).Post("/search", handlers.Search)

	// RESTful routing
	r.Route("/note", func(r chi.Router) {
		r.Use(middleware.URLFormat)
		r.With(requireAuthentication).Get("/create", handlers.ShowCreateNote)
		r.With(requireAuthentication).Post("/create", handlers.CreateNote)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(noteContext)
			r.Get("/", handlers.ShowNote)
			r.With(requireAuthentication).Get("/edit", handlers.ShowEditNote)
			r.With(requireAuthentication).Post("/", handlers.UpdateNote)
			r.With(requireAuthentication).Delete("/", handlers.DeleteNote)
		})
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/signup", handlers.ShowSignup)
		r.Post("/signup", handlers.Signup)
		r.Get("/login", handlers.ShowLogin)
		r.Post("/login", handlers.Login)
		r.With(requireAuthentication).Post("/logout", handlers.Logout)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return baseChain.Handler(r)
}
