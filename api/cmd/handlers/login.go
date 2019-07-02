package handlers

import (
	"database/sql"
	"net/http"

	"github.com/jongschneider/youtube-project/api/internal/platform/encryption"

	clientSVC "github.com/jongschneider/youtube-project/api/internal/platform/client"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/pkg/errors"
)

type loginResponse struct {
	web.Response
}

// Login lets a user login with a username and password
func Login(client *clientSVC.Client) http.HandlerFunc {
	log := client.Log()

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the content of the form
		err := r.ParseForm()
		if err != nil {
			log.WithError(err)
			web.RespondWithCodedError(w, r, http.StatusBadRequest, "bad request", err)
			return
		}

		email := r.FormValue("email")
		pass := r.FormValue("password")

		// Go out to the db and try to get the hashed password associated with the provided email
		hash, err := getPasswordByEmail(client, email)
		if err != nil {
			// The email was not in the db
			if err == sql.ErrNoRows {
				web.RespondWithCodedError(w, r, http.StatusBadRequest, "email does not exist", errors.Wrap(err, "login"))
				return
			}

			// Something else went wrong
			web.RespondWithCodedError(w, r, http.StatusInternalServerError, "", errors.Wrap(err, "login"))
			return
		}

		// Compare the hashed password we had in the db with a hashed version of the password the user provided.
		// If they are the same, we have a match!!!
		if !encryption.Compare(hash, pass) {
			web.RespondWithCodedError(w, r, http.StatusBadRequest, "invalid email/password", errors.Wrap(err, "login"))
			return
		}

		web.Respond(w, r, loginResponse{
			Response: web.Response{
				Message: "sucess",
			},
		}, http.StatusOK)
	}
}

// getPasswordByEmail gets the hashed password associated with the provided email
func getPasswordByEmail(client *clientSVC.Client, email string) (string, error) {
	db := client.DB()
	query := `SELECT password FROM users WHERE email = ?`

	target := []struct {
		Password string `db:"password"`
	}{}

	err := db.Select(&target, query, email)
	if err != nil {
		return "", err
	}
	if len(target) == 0 {
		return "", sql.ErrNoRows
	}

	return target[0].Password, nil
}
