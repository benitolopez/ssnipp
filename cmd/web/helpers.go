package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// serverError logs the detailed error message and stack trace, then sends a generic 500 Internal Server Error response to the user.
// If the application is in debug mode, it includes the error and stack trace in the response body.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	// Include the trace in the log entry.
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

	// If in debug mode, send the error and trace in the response body.
	if app.debug {
		body := fmt.Sprintf("%s\n%s", err, trace)
		http.Error(w, body, http.StatusInternalServerError)
		return
	}

	// And send a generic 500 Internal Server Error response.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific status code and corresponding description to the user.
// Used for sending responses like 400 "Bad Request" when there's a problem with the request.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// render renders a template, writing it to an internal buffer first to catch any errors.
// If there are no errors, it writes the buffered content to the http.ResponseWriter.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the template from the cache.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter.
	// If there's an error, call our serverError() helper and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If the template is written to the buffer without any errors, write the HTTP status code.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)
}

// newTemplateData returns a pointer to a templateData struct initialized with the current year,
// flash message, authentication status, signup allowance, and CSRF token.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		AllowSignup:     app.allowSignup,
		CSRFToken:       nosurf.Token(r),
	}
}

// decodePostForm parses the form data from the request and decodes it into the provided destination struct.
// It uses the form decoder to map the form values to the struct fields.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Parse the form data from the request.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Decode the form data into the destination struct.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If the target destination is invalid, raise a panic rather than returning the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, return them as normal.
		return err
	}

	return nil
}

// isAuthenticated checks if the current request is from an authenticated user by looking for a value
// in the request context. It returns true if the user is authenticated, otherwise false.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
