package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/hanna-yhchen/q-notes/internal/config"
)

var app *config.Application

type contextKey string

const (
	ContextKeyNote            = contextKey("note")
	ContextKeyIsAuthenticated = contextKey("isAuthenticated")
)

// NewHelpers sets the app configuration for the helpers package.
func NewHelpers(a *config.Application) {
	app = a
}

// serverError logs an error message and the stack trace for the current goroutine,
// and then reply a 500 internal server error to the client.
func ServerError(w http.ResponseWriter, err error) {
	stack := fmt.Sprintf("%s\n%s", err, debug.Stack())
	app.ErrorLog.Output(2, stack) // make the file name and line number reported one step back
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError replies a specific status code and corresponding description to
// the client.
func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound is a convenience wrapper around 404 client error.
func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}

// IsAuthenticated indicates whether the current request is from an authenticated
// user.
func IsAuthenticated(r *http.Request) bool {
	if isAuthenticated, ok := r.Context().Value(ContextKeyIsAuthenticated).(bool); ok {
		return isAuthenticated
	} else {
		return false
	}
}
