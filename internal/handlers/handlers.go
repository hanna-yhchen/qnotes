package handlers

import (
	"net/http"

	"github.com/hanna-yhchen/q-notes/internal/config"
	"github.com/hanna-yhchen/q-notes/internal/models"
	"github.com/hanna-yhchen/q-notes/internal/render"
)

var app *config.Application

// NewHandlers sets the app configuration for the handlers package.
func NewHandlers(a *config.Application) {
	app = a
}

func Home(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("Show Home Page")
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// GET /note/{id}
func ShowNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Show a specific note."))
	// render.Template(w, r, "note.page.tmpl", &models.TemplateData{})
}

// GET /note/create
func ShowCreateNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create new note."))
}

// POST /note/create
func CreateNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create new note."))
}

// PUT /note/{id}
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update a specific note."))
}

// DELETE /note/{id}
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("Call Delete")
	w.Write([]byte("Delete a specific note."))
}

// GET /user/signup
func ShowSignup(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "signup.page.tmpl", &models.TemplateData{})
}

// POST /user/signup
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create new note."))
}

// GET /user/login
func ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{})
}

// POST /user/login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Log in."))
}

// POST /user/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Log out."))
}
