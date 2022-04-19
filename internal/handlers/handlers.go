package handlers

import (
	"errors"
	"net/http"

	"github.com/hanna-yhchen/q-notes/internal/config"
	"github.com/hanna-yhchen/q-notes/internal/forms"
	"github.com/hanna-yhchen/q-notes/internal/helpers"
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

// GET /note/{id}/edit
func ShowEditNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Edit a specific note."))
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
	render.Template(w, r, "signup.page.tmpl", &models.TemplateData{Form: forms.New(nil)})
}

// POST /user/signup
func Signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MinLength("password", 6)
	form.MatchPattern("email", forms.EmailRegex)

	// If the form is not valid, re-render the page with client's input.
	if !form.IsValid() {
		render.Template(w, r, "signup.page.tmpl", &models.TemplateData{Form: form})
		return
	}

	if err := app.UserModel.Insert(form.Get("name"), form.Get("email"), form.Get("password")); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "The email address is already in use.")
			render.Template(w, r, "signup.page.tmpl", &models.TemplateData{Form: form})
		} else {
			helpers.ServerError(w, err)
		}
		return
	}

	app.Session.Put(r, "flash", "You have successfully signed up! Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// GET /user/login
func ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: forms.New(nil)})
}

// POST /user/login
func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.UserModel.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Incorrect email or password.")
			render.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: form})
		} else {
			helpers.ServerError(w, err)
		}
		return
	}

	app.Session.Put(r, "authenticatedUserID", id)
	app.Session.Put(r, "flash", "You've been logged in successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /user/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Remove(r, "authenticatedUserID")
	app.Session.Put(r, "flash", "You've been logged out.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
