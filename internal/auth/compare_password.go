package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

// Check a password string against a hashed password
func CheckPasswordHash(password, hashed_password string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hashed_password)
	if err != nil {
		err := fmt.Errorf("couldn't check password hash: %s", err)
		return false, err
	}
	return match, nil
}
