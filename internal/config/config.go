package config

import (
	"html/template"
	"log"

	"github.com/golangcollege/sessions"
)

// Application holds the app-wide dependencies.
type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Session       *sessions.Session
	TemplateCache map[string]*template.Template
	// noteModel
	// userModel
}
