package main

import (
	"net/http"
	"net/url"
	"testing"

	"ssnipp.com/internal/assert"
)

// TestPing tests the /ping endpoint to ensure it returns a 200 OK status and "OK" body.
func TestPing(t *testing.T) {
	// Create a new instance of the application struct containing mocked dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET request to the /ping endpoint.
	code, _, body := ts.get(t, "/ping")

	// Assert that the status code is 200 OK.
	assert.Equal(t, code, http.StatusOK)

	// Assert that the response body is "OK".
	assert.Equal(t, body, "OK")
}

// TestSnippetView tests the /view/{id} endpoint with various IDs to check for proper handling.
func TestSnippetView(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our application for different URLs.
	tests := []struct {
		name     string // Name of the test case.
		urlPath  string // URL path to test.
		wantCode int    // Expected HTTP status code.
		wantBody string // Expected response body (if any).
	}{
		{
			name:     "Valid ID",
			urlPath:  "/view/1",
			wantCode: http.StatusOK,
			wantBody: "console.log();",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/view/",
			wantCode: http.StatusNotFound,
		},
	}

	// Iterate over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a GET request to the test URL path.
			code, _, body := ts.get(t, tt.urlPath)

			// Assert that the status code matches the expected value.
			assert.Equal(t, code, tt.wantCode)

			// If an expected body is provided, assert that it is contained in the response body.
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

// TestUserSignup tests the /signup endpoint with various form submissions to check for proper handling.
func TestUserSignup(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET request to the /signup endpoint to retrieve a valid CSRF token.
	_, _, body := ts.get(t, "/signup")
	validCSRFToken := extractCSRFToken(t, body)

	// Define valid test data constants.
	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/signup' method='POST' novalidate>"
	)

	// Set up some table-driven tests to check the responses for different form submissions.
	tests := []struct {
		name         string // Name of the test case.
		userName     string // User name to test.
		userEmail    string // User email to test.
		userPassword string // User password to test.
		csrfToken    string // CSRF token to use.
		wantCode     int    // Expected HTTP status code.
		wantFormTag  string // Expected form tag in the response body (if any).
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	// Iterate over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a form with the test data.
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			// Make a POST request to the /signup endpoint with the form data.
			code, _, body := ts.postForm(t, "/signup", form)

			// Assert that the status code matches the expected value.
			assert.Equal(t, code, tt.wantCode)

			// If an expected form tag is provided, assert that it is contained in the response body.
			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}
