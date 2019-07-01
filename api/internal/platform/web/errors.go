package web

import (
	"net/http"

	"github.com/go-chi/render"
)

type errResponse struct {
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	Status    int    `json:"status,omitempty"`
}

// RespondWithCodedError responds with ambiguous error JSON
func RespondWithCodedError(w http.ResponseWriter, r *http.Request, httpStatus int, code string, err error) {
	serr := http.StatusText(httpStatus)

	if ct := w.Header().Get("Content-Type"); ct == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	w.WriteHeader(httpStatus)
	render.JSON(w, r, errResponse{
		Error:     serr,
		Status:    httpStatus,
		ErrorCode: code,
	})
}
