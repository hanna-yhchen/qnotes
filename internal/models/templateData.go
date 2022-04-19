package models

type TemplateData struct {
	CSRFToken string
	Flash     string
	// Form
	// IsAuthenticated bool
	Note  Note
	Notes []Note
}
