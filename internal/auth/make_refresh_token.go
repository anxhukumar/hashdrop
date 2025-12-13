package auth

import (
	"crypto/rand"
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
