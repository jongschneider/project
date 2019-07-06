package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jongschneider/youtube-project/api/internal/platform/user"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/pkg/errors"
)

type updateResponse struct {
	web.Response
}

// Update updates a user with a username and password
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse the content of the url
	id := chi.URLParam(r, "ID")

	userID, err := strconv.Atoi(id)
	if err != nil {
		h.log.WithError(err)
		web.RespondWithCodedError(w, r, http.StatusBadRequest, "bad request", err)
		return
	}

	// Parse the content of the form
	target := user.User{}
	err = json.NewDecoder(r.Body).Decode(&target)
	defer r.Body.Close()
	if err != nil {
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusBadRequest, "malformed request", errors.Wrap(err, "decode"))
		return
	}

	target.ID = userID
	// Go out to the db and try to get the hashed password associated with the provided email
	err = user.Update(h.db, target)
	if err != nil {
		// Something else went wrong
		h.log.WithError(err).Info()
		web.RespondWithCodedError(w, r, http.StatusInternalServerError, "", errors.Wrap(err, "update"))
		return
	}

	web.Respond(w, r, updateResponse{
		Response: web.Response{
			Message: "success",
		},
	}, http.StatusNoContent)
}
