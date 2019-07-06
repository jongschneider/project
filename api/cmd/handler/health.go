package handler

import (
	"net/http"

	"github.com/jongschneider/youtube-project/api/internal/platform/web"
)

// Health is the health check for the application
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if err := h.db.PingContext(r.Context()); err != nil {
		web.RespondWithCodedError(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err)
		return
	}

	response := web.Response{
		Message: "Healthy",
	}

	web.Respond(w, r, response, http.StatusOK)
	return
}
