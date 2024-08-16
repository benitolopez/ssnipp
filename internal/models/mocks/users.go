package mocks

import (
	"ssnipp.com/internal/models"
)

type UserModel struct{}

// Insert is a mock implementation of the Insert method. It returns an ErrDuplicateEmail
// error if the email is "dupe@example.com". Otherwise, it returns nil.
func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

// Authenticate is a mock implementation of the Authenticate method. It returns
// user ID 1 and nil error if the email is "alice@example.com" and the password
// is "pa$$word". Otherwise, it returns 0 and an ErrInvalidCredentials error.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

// Exists is a mock implementation of the Exists method. It returns true and nil
// error if the user ID is 1. Otherwise, it returns false and nil error.
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
