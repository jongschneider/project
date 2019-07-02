package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// New creates a new JWT token
func New(key string, id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = id

	token.Claims = claims

	return token.SignedString([]byte(key))

}
