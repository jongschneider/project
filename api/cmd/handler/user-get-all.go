package handler

import (
	"database/sql"
	"net/http"

	"github.com/jongschneider/youtube-project/api/internal/platform/user"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/pkg/errors"
)

type getAllResponse struct {
	web.Response
	Users []user.User `json:"users,omitempty"`
}

// GetAllUsers gets all users
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Go out to the db and try to get the hashed password associated with the provided email
	users, err := user.GetAll(h.db)
	if err != nil {
		// The user was not in the db
		if err == sql.ErrNoRows {
			h.log.WithError(err).Info()
			web.RespondWithCodedError(w, r, http.StatusBadRequest, "users do not exist", errors.Wrap(err, "get all users"))
			return
		}

		// Something else went wrong
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusInternalServerError, "", errors.Wrap(err, "get all users"))
		return
	}

	web.Respond(w, r, getAllResponse{
		Response: web.Response{
			Message: "success",
		},
		Users: users,
	}, http.StatusOK)
}
