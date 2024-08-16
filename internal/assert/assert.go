package assert

import (
	"strings"
	"testing"
)

// Equal is a generic test helper function that checks if the actual value is equal to the expected value.
// If they are not equal, it reports an error with the actual and expected values.
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// StringContains checks if the actual string contains the expected substring.
// If it does not, it reports an error with the actual string and the expected substring.
func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

// NilError checks if the actual error is nil.
// If it is not nil, it reports an error with the actual error.
func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
