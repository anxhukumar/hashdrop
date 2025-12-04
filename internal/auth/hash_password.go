package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

// Hash the password
func HashedPassword(password string) (string, error) {
	if len(password) == 0 {
		err := fmt.Errorf("password field is empty")
		return "", err
	}
	if len(password) < 8 {
		err := fmt.Errorf("length of password is less than 8 characters")
		return "", err
	}
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		err := fmt.Errorf("couldn't hash password: %s", err)
		return "", err
	}
	return hash, nil
}
