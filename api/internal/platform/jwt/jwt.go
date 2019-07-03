package jwt

import (
	"errors"
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

func ParseToken(key, val string) (int, error) {
	token, err := jwt.Parse(val, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return -1, err
	}

	if !token.Valid {
		return -1, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return -1, errors.New("Invalid claims")
	}

	id := int(claims["id"].(float64))

	return id, nil
}
