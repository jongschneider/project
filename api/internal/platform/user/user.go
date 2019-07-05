package user

import (
	"database/sql"

	"github.com/jongschneider/youtube-project/api/internal/platform/database"
)

// User represents a user in the DB
type User struct {
	ID       int    `db:"id"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

// GetUserByEmail gets a user associated with the provided email
func GetUserByEmail(db *database.DB, email string) (User, error) {
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

// GetUserByID gets gets a user associated with the provided id
func GetUserByID(db *database.DB, id int) (User, error) {
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
