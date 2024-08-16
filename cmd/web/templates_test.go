package main

import (
	"testing"
	"time"

	"ssnipp.com/internal/assert"
)

// TestHumanDate tests the humanDate function to ensure it returns the correct formatted string
// representation of a time.Time object.
func TestHumanDate(t *testing.T) {
	// Create a slice of anonymous structs containing the test case name,
	// input to our humanDate() function (the tm field), and expected output
	// (the want field).
	tests := []struct {
		name string    // Name of the test case.
		tm   time.Time // Input time to test.
		want string    // Expected output string.
	}{
		{
			name: "UTC", // Test case for a specific UTC time.
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty", // Test case for an empty time value.
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET", // Test case for a specific CET time.
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for each test case. The
		// first parameter to this is the name of the test (which is used to
		// identify the sub-test in any log output) and the second parameter is
		// an anonymous function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			// Call the humanDate function with the test case input.
			hd := humanDate(tt.tm)

			// Assert that the result matches the expected output.
			assert.Equal(t, hd, tt.want)
		})
	}
}

// TestHumanDate tests the getLanguageLabel function to ensure it returns the correct language label
func TestGetLanguageLabel(t *testing.T) {
	// Create a slice of anonymous structs containing the test case name,
	// input to our humanDate() function (the tm field), and expected output
	// (the want field).
	tests := []struct {
		name string // Name of the test case.
		lang string // Input language to test.
		want string // Expected output string.
	}{
		{
			name: "Valid language key",
			lang: "javascript",
			want: "JavaScript",
		},
		{
			name: "Invalid language key",
			lang: "latin",
			want: "Plain Text",
		},
		{
			name: "Empty language key",
			lang: "",
			want: "Plain Text",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for each test case. The
		// first parameter to this is the name of the test (which is used to
		// identify the sub-test in any log output) and the second parameter is
		// an anonymous function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			// Call the getLanguageLabel function with the test case input.
			label := getLanguageLabel(tt.lang)

			// Assert that the result matches the expected output.
			assert.Equal(t, label, tt.want)
		})
	}
}
