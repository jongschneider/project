package handler

import (
	"database/sql"
	"net/http"

	"github.com/jongschneider/youtube-project/api/internal/platform/encryption"
	"github.com/jongschneider/youtube-project/api/internal/platform/user"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/pkg/errors"
)

type loginResponse struct {
	web.Response
}

// Login lets a user login with a username and password
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the content of the form
	err := r.ParseForm()
	if err != nil {
		h.log.WithError(err)
		web.RespondWithCodedError(w, r, http.StatusBadRequest, "bad request", err)
		return
	}

	email := r.FormValue("email")
	pass := r.FormValue("password")

	// Go out to the db and try to get the hashed password associated with the provided email
	u, err := user.GetByEmail(h.db, email)
	if err != nil {
		// The email was not in the db
		if err == sql.ErrNoRows {
			h.log.WithError(err).Info()
			web.RespondWithCodedError(w, r, http.StatusBadRequest, "email does not exist", errors.Wrap(err, "login"))
			return
		}

		// Something else went wrong
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusInternalServerError, "", errors.Wrap(err, "login"))
		return
	}

	// Compare the hashed password we had in the db with a hashed version of the password the user provided.
	// If they are the same, we have a match!!!
	if !encryption.Compare(u.Password, pass) {
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusBadRequest, "invalid email/password", errors.Wrap(err, "login"))
		return
	}

	// token, err := jwtSVC.New(h.key, u.ID)
	// if err != nil {
	// 	h.log.WithError(err).Info()
	// 	web.RespondWithCodedError(w, r, http.StatusInternalServerError, "could not issue JWT token", errors.Wrap(err, "login"))
	// 	return
	// }

	web.Respond(w, r, loginResponse{
		Response: web.Response{
			Message: "success",
		},
		// Token: token,
	}, http.StatusOK)
}
