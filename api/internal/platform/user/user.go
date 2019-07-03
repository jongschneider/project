package user

import (
	"database/sql"

	clientSVC "github.com/jongschneider/youtube-project/api/internal/platform/client"
)

type User struct {
	ID       int    `db:"id"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

// GetUserByEmail gets a user associated with the provided email
func GetUserByEmail(client *clientSVC.Client, email string) (User, error) {
	db := client.DB()
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
func GetUserByID(client *clientSVC.Client, id int) (User, error) {
	db := client.DB()
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
