package encryption

import (
	"crypto/rand"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"golang.org/x/crypto/argon2"
)

// Generate DEK from passphrase for no-vault users
func GenerateDEKfromPassphrase(passphrase string) (key []byte, salt []byte, err error) {
	passphrase_bytes := []byte(passphrase)

	salt = make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		return nil, nil, err
	}

	key = argon2.IDKey(
		passphrase_bytes,
		salt,
		config.ArgonTime,
		config.ArgonMemory,
		config.ArgonThreads,
		config.ArgonKeyLen,
	)

	return key, salt, nil
}

// Generate random DEK for vault users
func GenerateRandomDEK() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, err
}
