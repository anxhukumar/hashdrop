package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {

	// Generate random key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("error generating random key for refresh token: %w", err)
	}

	keyStr := hex.EncodeToString(key)
	return keyStr, nil
}

// HashRefreshToken hashes token using a secret key
func HashRefreshToken(token string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
