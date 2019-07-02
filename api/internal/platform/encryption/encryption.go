package encryption

import "golang.org/x/crypto/bcrypt"

const salt = 10

func Encrypt(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), salt)
	return string(hash), err
}

func Compare(hash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return false
	}

	return true
}
