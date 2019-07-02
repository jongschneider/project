package handlers

import (
	"net/http"

	clientSVC "github.com/jongschneider/youtube-project/api/internal/platform/client"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
)

// Health is the health check for the application
func Health(client *clientSVC.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := client.DB().PingContext(r.Context()); err != nil {
			web.RespondWithCodedError(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err)
			return
		}

		web.Respond(w, r, "Healthy", http.StatusOK)
		return
	}
}
