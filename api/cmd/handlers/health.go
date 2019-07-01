package handlers

import (
	"net/http"

	"github.com/jongschneider/youtube-project/api/internal/platform/web"
)

// Health is the health check for the application
func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		web.Respond(w, r, "Healthy", http.StatusOK)
	}
}
