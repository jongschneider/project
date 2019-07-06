package user

import (
	"database/sql"

	"github.com/jongschneider/youtube-project/api/internal/platform/database"
	"github.com/jongschneider/youtube-project/api/internal/platform/encryption"
	"github.com/pkg/errors"
)

// User represents a user in the DB
type User struct {
	ID       int    `json:"id" db:"id"`
	Password string `json:"-" db:"password"`
	Email    string `json:"email" db:"email"`
}

// GetByEmail gets a user associated with the provided email
func GetByEmail(db *database.DB, email string) (User, error) {
	query := `SELECT id, password, email FROM users WHERE email = ?`

	target := []User{}

	err := db.Select(&target, query, email)
	if err != nil {
		return User{}, err
	}

	if len(target) == 0 {
		return User{}, sql.ErrNoRows
	}

	return target[0], nil
}

// GetByID gets gets a user associated with the provided id
func GetByID(db *database.DB, id int) (User, error) {
	query := `SELECT id, password, email FROM users WHERE id = ?`

	target := []User{}

	err := db.Select(&target, query, id)
	if err != nil {
		return User{}, err
	}

	if len(target) == 0 {
		return User{}, sql.ErrNoRows
	}

	return target[0], nil
}

// GetAll gets gets a user associated with the provided id
func GetAll(db *database.DB) ([]User, error) {
	query := `SELECT id, password, email FROM users`

	target := []User{}

	err := db.Select(&target, query)
	if err != nil {
		return target, err
	}

	if len(target) == 0 {
		return target, sql.ErrNoRows
	}

	return target, nil
}

// Delete deletes a user from the user table
func Delete(db *database.DB, id int) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert creates a new user.
func Insert(db *database.DB, u User) error {
	query := `INSERT INTO users (email, password) VALUE ( ?, ? )`

	hash, err := encryption.Encrypt(u.Password)
	if err != nil {
		return errors.Wrap(err, "insert")
	}
	_, err = db.Exec(query, u.Email, hash)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing user.
func Update(db *database.DB, u User) error {
	query := `UPDATE users SET email = ?, password = ? WHERE id = ?`

	hash, err := encryption.Encrypt(u.Password)
	if err != nil {
		return errors.Wrap(err, "update")
	}
	_, err = db.Exec(query, u.Email, hash, u.ID)
	if err != nil {
		return err
	}

	return nil
}
