package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// commonHeaders middleware sets various security-related HTTP headers on the response.
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Content Security Policy (CSP) header.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self'; font-src 'self'; img-src 'self' data:;")

		// Set Referrer Policy header.
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// Set other security headers.
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// logRequest middleware logs details of each incoming HTTP request.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		// Log request details.
		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// recoverPanic middleware recovers from any panics and returns a 500 Internal Server Error.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Recover from panic, if any.
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")

				// Return a 500 Internal Server Error response.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// requireAuthentication middleware checks if a user is authenticated, otherwise redirects to login page.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			// Store the path the user is trying to access.
			app.sessionManager.Put(r.Context(), "redirectPathAfterLogin", r.URL.Path)

			// Redirect to login page.
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Set Cache-Control header to prevent caching of authenticated pages.
		w.Header().Add("Cache-Control", "no-store")

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// noSurf middleware sets up CSRF protection using the nosurf package.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// Set custom CSRF cookie attributes.
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
	})

	return csrfHandler
}

// authenticate middleware checks if a user is authenticated and adds authentication information to the request context.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticatedUserID value from the session.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			// If no authenticatedUserID, call the next handler in the chain.
			next.ServeHTTP(w, r)
			return
		}

		// Check if a user with that ID exists in the database.
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// If a matching user is found, add authentication information to the request context.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
