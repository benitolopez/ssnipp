package models

import (
	"testing"

	"ssnipp.com/internal/assert"
)

// TestUserModelExists tests the Exists method of the UserModel.
func TestUserModelExists(t *testing.T) {
	// Skip the test if the "-short" flag is provided when running the test.
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID", // Test case with a valid user ID.
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID", // Test case with a zero user ID.
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID", // Test case with a non-existent user ID.
			userID: 2,
			want:   false,
		},
	}

	// Iterate over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() helper function to get a connection pool to
			// our test database.
			db := newTestDB(t)

			// Create a new instance of the UserModel.
			m := UserModel{db}

			// Call the UserModel.Exists() method and check that the return
			// value and error match the expected values for the sub-test.
			exists, err := m.Exists(tt.userID)

			// Assert that the result matches the expected value.
			assert.Equal(t, exists, tt.want)

			// Assert that there is no error.
			assert.NilError(t, err)
		})
	}
}
