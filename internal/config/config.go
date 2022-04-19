package config

import (
	"html/template"
	"log"

	"github.com/golangcollege/sessions"
	"github.com/hanna-yhchen/q-notes/internal/models/mysql"
)

// Application holds the app-wide dependencies.
type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Session       *sessions.Session
	TemplateCache map[string]*template.Template
	NoteModel     *mysql.NoteModel
	UserModel     *mysql.UserModel
}
