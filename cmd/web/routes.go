package main

import (
	"net/http"

	"github.com/justinas/alice"
	"ssnipp.com/ui"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create an HTTP handler which
	// serves the embedded files in ui.Files. Our static files are contained
	// in the "static" folder of the ui.Files embedded filesystem. This means
	// that any requests that start with /static/ can be passed directly to the
	// file server, and the corresponding static file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Add a route for the ping handler.
	mux.HandleFunc("GET /ping", ping)

	// Create a middleware chain for dynamic routes which includes the session manager,
	// CSRF protection, and authentication middleware.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Add routes for viewing snippets and user login.
	mux.Handle("GET /view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /login", dynamic.ThenFunc(app.userLoginPost))

	// If the allowSignup configuration setting is true, add routes for user signup.
	// Otherwise, these routes will not be available.
	if app.allowSignup {
		mux.Handle("GET /signup", dynamic.ThenFunc(app.userSignup))
		mux.Handle("POST /signup", dynamic.ThenFunc(app.userSignupPost))
	}

	// Create a new middleware chain for protected (authenticated-only) routes
	// which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	// Add routes for home, snippet creation, and user logout.
	mux.Handle("GET /{$}", protected.ThenFunc(app.home))
	mux.Handle("POST /create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /logout", protected.ThenFunc(app.userLogoutPost))

	// Create a standard middleware chain which includes the panic recovery,
	// request logging, and common security headers middleware.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the servemux wrapped with the standard middleware chain.
	return standard.Then(mux)
}
