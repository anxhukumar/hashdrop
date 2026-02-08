package otp

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

func HashOTP(otp string, secret string) string {
	sum := sha256.Sum256([]byte(otp + secret))
	return hex.EncodeToString(sum[:])
}

func VerifyOTP(inputOTP, storedHash, secret string) bool {
	inputHash := HashOTP(inputOTP, secret)

	// Constant-time compare
	return subtle.ConstantTimeCompare([]byte(inputHash), []byte(storedHash)) == 1
}
