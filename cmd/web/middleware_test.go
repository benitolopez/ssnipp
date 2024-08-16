package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"ssnipp.com/internal/assert"
)

// TestCommonHeaders tests the commonHeaders middleware to ensure it sets the correct headers
// and properly calls the next handler in the chain.
func TestCommonHeaders(t *testing.T) {
	// Initialize a new httptest.ResponseRecorder to record the response.
	rr := httptest.NewRecorder()

	// Create a dummy HTTP GET request.
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our commonHeaders middleware.
	// This handler simply writes a 200 status code and an "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock HTTP handler to our commonHeaders middleware.
	// Since commonHeaders returns a http.Handler, we can call its ServeHTTP() method,
	// passing in the http.ResponseRecorder and dummy http.Request to execute it.
	commonHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method on the http.ResponseRecorder to get the results of the test.
	rs := rr.Result()

	// Check that the middleware has correctly set the Content-Security-Policy header on the response.
	expectedValue := "default-src 'self'; style-src 'self'; font-src 'self'; img-src 'self' data:;"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// Check that the middleware has correctly set the Referrer-Policy header on the response.
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// Check that the middleware has correctly set the X-Content-Type-Options header on the response.
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// Check that the middleware has correctly set the X-Frame-Options header on the response.
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	// Check that the middleware has correctly set the X-XSS-Protection header on the response.
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	// Check that the middleware has correctly set the Server header on the response.
	expectedValue = "Go"
	assert.Equal(t, rs.Header.Get("Server"), expectedValue)

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// Read and trim the response body.
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	// Check that the response body is "OK".
	assert.Equal(t, string(body), "OK")
}
