package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const hashBytesLimit = 16 // 128 bits. Safe and reasonable.

// Generate hash from UserID to use it as prefix of s3key
func GenerateUserIDHash(userID string, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(userID))
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum[:hashBytesLimit])
}
