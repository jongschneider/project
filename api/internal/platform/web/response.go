package web

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Message string `json:"message,omitempty"`
}

// Respond converts a Go value to JSON and sends it to the client.
func Respond(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) error {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	render.JSON(w, r, data)

	return nil
}
