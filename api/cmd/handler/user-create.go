package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jongschneider/youtube-project/api/internal/platform/user"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/pkg/errors"
)

type createResponse struct {
	web.Response
}

// Create create a user with a username and password
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse the content of the form
	var target user.User
	err := json.NewDecoder(r.Body).Decode(&target)
	defer r.Body.Close()
	if err != nil {
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusBadRequest, "malformed request", errors.Wrap(err, "decode"))
		return
	}

	// Go out to the db and try to get the hashed password associated with the provided email
	err = user.Insert(h.db, target)
	if err != nil {
		// Something else went wrong
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusInternalServerError, "", errors.Wrap(err, "create"))
		return
	}

	web.Respond(w, r, createResponse{
		Response: web.Response{
			Message: "success",
		},
	}, http.StatusNoContent)
}
