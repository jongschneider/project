package handlers

import (
	"net/http"

	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/render"
)

type loginResponse struct {
}

// Login lets a user login with a username and password
func Login(log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.WithError(err)
			web.RespondWithCodedError(w, r, http.StatusBadRequest, "bad request", err)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		log.Println("email: ", email)
		log.Println("password: ", password)
		render.Respond(w, r, "blah")
	}
}
