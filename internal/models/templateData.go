package models

import "github.com/hanna-yhchen/q-notes/internal/forms"

type TemplateData struct {
	CSRFToken       string
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	Note            *Note
	Notes           []*Note
	NoteID          int
	IsNew           bool
}
