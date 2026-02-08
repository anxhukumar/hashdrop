package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Generates a 6-digit numeric otp as string
func GenerateOTP() (string, error) {
	// 0 to 999999
	max := big.NewInt(1000000)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Zero-pad to 6 digits
	return fmt.Sprintf("%06d", n.Int64()), nil
}
