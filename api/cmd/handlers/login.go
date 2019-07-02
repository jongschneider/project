package handlers

import (
	"net/http"

	clientSVC "github.com/jongschneider/youtube-project/api/internal/platform/client"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"

	"github.com/go-chi/render"
)

type loginResponse struct {
}

// Login lets a user login with a username and password
func Login(client *clientSVC.Client) http.HandlerFunc {
	log := client.Log()

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
