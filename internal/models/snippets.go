package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(content string, language string) (int, error)
	Get(id int) (Snippet, error)
}

// Snippet represents a single code snippet. The fields correspond to the columns
// in our MySQL snippets table.
type Snippet struct {
	ID       int
	Content  string
	Created  time.Time
	Language string
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert adds a new snippet to the database and returns the ID of the newly inserted record.
func (m *SnippetModel) Insert(content string, language string) (int, error) {
	// SQL statement to insert a new snippet into the database.
	stmt := `INSERT INTO snippets (content, created, language)
    VALUES(?, UTC_TIMESTAMP(), ?)`

	// Execute the SQL statement using the Exec() method. The content and language
	// parameters will be substituted into the placeholders in the SQL statement.
	result, err := m.DB.Exec(stmt, content, language)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted record.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Convert the ID from int64 to int and return it.
	return int(id), nil
}

// Get retrieves a specific snippet based on its ID.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	// SQL statement to retrieve a snippet by its ID.
	stmt := `SELECT id, content, created, language FROM snippets
    WHERE id = ?`

	// Execute the SQL statement using the QueryRow() method, passing in the ID
	// as the value for the placeholder parameter. This returns a pointer to a sql.Row object.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a new zeroed Snippet struct.
	var s Snippet

	// Copy the values from the sql.Row object to the Snippet struct using the Scan() method.
	err := row.Scan(&s.ID, &s.Content, &s.Created, &s.Language)
	if err != nil {
		// If the query returns no rows, return a ErrNoRecord error.
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	// Return the filled Snippet struct.
	return s, nil
}
