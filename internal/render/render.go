package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/hanna-yhchen/q-notes/internal/config"
	"github.com/hanna-yhchen/q-notes/internal/helpers"
	"github.com/hanna-yhchen/q-notes/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.Application

// NewRenderer sets the app configuration for the render package.
func NewRenderer(a *config.Application) {
	app = a
}

// Template renders the template from cache with the given filename and template data.
func Template(w http.ResponseWriter, r *http.Request, name string, td *models.TemplateData) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		helpers.ServerError(w, fmt.Errorf("the template with filename %q does not exist", name))
		return
	}

	// Render templates in two stages to avoid unexpected behaviors such like
	// showing half-completed page.
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, addDefaultData(td, r))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	buf.WriteTo(w)
}

// NewTemplateCache creates a map from template name to template cache.
func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

// addDefaultData adds common dynamic data to the template data.
func addDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	if td == nil {
		td = &models.TemplateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = helpers.IsAuthenticated(r)
	td.Flash = app.Session.PopString(r, "flash")
	return td
}
