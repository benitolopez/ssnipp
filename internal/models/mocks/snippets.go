package mocks

import (
	"time"

	"ssnipp.com/internal/models"
)

// mockSnippet is a sample Snippet used for mocking purposes in tests.
var mockSnippet = models.Snippet{
	ID:       1,
	Content:  "console.log();",
	Created:  time.Now(),
	Language: "javascript",
}

// SnippetModel is a mock implementation of the SnippetModel interface.
type SnippetModel struct{}

// Insert is a mock implementation of the Insert method. It returns a fixed ID and nil error.
func (m *SnippetModel) Insert(content string, language string) (int, error) {
	return 2, nil
}

// Get is a mock implementation of the Get method. It returns the mockSnippet if the ID is 1,
// otherwise it returns an empty Snippet and an ErrNoRecord error.
func (m *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}
