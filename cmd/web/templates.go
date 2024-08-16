package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"ssnipp.com/internal/models"
	"ssnipp.com/ui"
)

// templateData type acts as the holding structure for any dynamic data that
// we want to pass to our HTML templates. It contains fields for the current year,
// snippet data, form data, flash messages, authentication status, CSRF token,
// signup allowance, and available languages.
type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	AllowSignup     bool
	Languages       []Language
}

// newTemplateCache creates a template cache as a map. The map's keys are the names of the templates
// (e.g., 'home.html') and the values are the parsed template sets (html/template.Template).
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded filesystem
	// which match the pattern 'html/pages/*.html'. This gives us a slice of all the 'page'
	// templates for the application.
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through the page templates.
	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we want to parse.
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		// Use ParseFS() to parse the template files from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the parsed template set to the cache using the name of the page template as the key.
		cache[name] = ts
	}

	return cache, nil
}

// humanDate function returns a nicely formatted string representation of a time.Time object.
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}

	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is essentially
// a string-keyed map which acts as a lookup between the names of our custom template functions
// and the functions themselves.
var functions = template.FuncMap{
	"humanDate":        humanDate,
	"getLanguageLabel": getLanguageLabel,
}
