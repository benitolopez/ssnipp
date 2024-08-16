package models

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// newTestDB creates a new database connection for testing purposes.
func newTestDB(t *testing.T) *sql.DB {
	// Load the .env file into the environment...
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
		os.Exit(1)
	}

	// Get the database username from the environment variable.
	dbTestUser := os.Getenv("DB_TEST_USERNAME")
	if dbTestUser == "" {
		t.Fatal("DB_TEST_USERNAME environment variable not set")
	}

	// Get the database password from the environment variable.
	dbTestPassword := os.Getenv("DB_TEST_PASSWORD")
	if dbTestPassword == "" {
		t.Fatal("DB_TEST_PASSWORD environment variable not set")
		os.Exit(1)
	}

	// Get the database name from the environment variable.
	dbTestDatabase := os.Getenv("DB_TEST_DATABASE")
	if dbTestDatabase == "" {
		t.Fatal("DB_TEST_DATABASE environment variable not set")
		os.Exit(1)
	}

	// Construct the DSN for the MySQL connection.
	dsn := dbTestUser + ":" + dbTestPassword + "@/" + dbTestDatabase + "?parseTime=true&multiStatements=true"

	// Open a connection to the database.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from the file and execute the statements.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	// Register a cleanup function which will be called when the test finishes.
	// This function reads and executes the teardown script and closes the database connection pool.
	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	// Return the database connection pool.
	return db
}
