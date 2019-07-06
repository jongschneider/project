package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jongschneider/youtube-project/api/internal/platform/user"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/pkg/errors"
)

type deleteResponse struct {
	web.Response
}

// Delete deletes a user with a username and password
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse the content of the url
	id := chi.URLParam(r, "ID")
	userID, err := strconv.Atoi(id)
	if err != nil {
		h.log.WithError(err)
		web.RespondWithCodedError(w, r, http.StatusBadRequest, "bad request", err)
		return
	}
	// Go out to the db and try to get the hashed password associated with the provided email
	u, err := user.GetByID(h.db, userID)

	if err != nil {
		// The user was not in the db
		if err == sql.ErrNoRows {
			h.log.WithError(err).Info()
			web.RespondWithCodedError(w, r, http.StatusBadRequest, "user does not exist", errors.Wrap(err, "get user"))
			return
		}

		// Something else went wrong
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusInternalServerError, "", errors.Wrap(err, "login"))
		return
	}

	err = user.Delete(h.db, u.ID)
	if err != nil {
		// Something else went wrong
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusInternalServerError, "", errors.Wrap(err, "delete"))
		return
	}

	web.Respond(w, r, deleteResponse{
		Response: web.Response{
			Message: "success",
		},
	}, http.StatusNoContent)
}
