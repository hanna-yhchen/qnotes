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

	r.Use(app.Session.Enable, noSurf)

	r.Get("/", handlers.Home)

	// RESTful routing
	r.Route("/note", func(r chi.Router) {
		r.Use(middleware.URLFormat)
		//r.With().Get()
		r.Get("/create", handlers.ShowCreateNote) // Need Auth
		r.Post("/create", handlers.CreateNote)    // Need Auth

		r.Route("/{id}", func(r chi.Router) {
			r.Use(noteContext)
			r.Get("/", handlers.ShowNote)
			r.Put("/", handlers.UpdateNote)    // Need Auth
			r.Delete("/", handlers.DeleteNote) // Need Auth
		})
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/signup", handlers.ShowSignup)
		r.Post("/signup", handlers.Signup)
		r.Get("/login", handlers.ShowLogin)
		r.Post("/login", handlers.Login)
		r.Post("/logout", handlers.Logout) // Need Auth
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return baseChain.Handler(r)
}
