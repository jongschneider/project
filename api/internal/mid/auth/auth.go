package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jongschneider/youtube-project/api/internal/platform/database"
	jwtSVC "github.com/jongschneider/youtube-project/api/internal/platform/jwt"
	"github.com/jongschneider/youtube-project/api/internal/platform/user"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
)

// ctxKey represents the type of value for the context key.
type authCtxKey string

// Key is used to store/retrieve a Claims value from a context.Context.
const tokenKey authCtxKey = "token"

// JWTMiddleware is middleware that validates a JWT token found in the header of a request.
func JWTMiddleware(key string, db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Parse the X-App-Token  header. Expected header is of
			// the format `Bearer <token>`.
			parts := strings.Split(r.Header.Get("X-App-Token"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Println("AUTH ERROR: x-app-token")
				err := errors.New("expected X-App-Token header format: Bearer <token>")
				web.RespondWithCodedError(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
				return
			}

			id, err := jwtSVC.ParseToken(key, parts[1])
			if err != nil {
				log.Println("AUTH ERROR: jwtSVC.ParseToken")
				web.RespondWithCodedError(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
				return
			}

			_, err = user.GetUserByID(db, id)
			if err != nil {
				log.Println("AUTH ERROR: user.GetUserByID")
				web.RespondWithCodedError(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
				return
			}

			// Add claims to the context so they can be retrieved later.
			ctx := context.WithValue(r.Context(), tokenKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
