package handlers

import (
	"fmt"
	"time"

	"github.com/anxhukumar/hashdrop/internal/auth"
	"github.com/anxhukumar/hashdrop/internal/database"
)

func GetJWTToken(userData database.User, jwtSecret string, tokenExpiry string) (string, error) {

	// Convert token expiry duration to appropirate format
	expiry, err := time.ParseDuration(tokenExpiry)
	if err != nil {
		return "", fmt.Errorf("error parsing access token expiry duration string")
	}

	// Fetch access token
	token, err := auth.MakeJWT(
		userData.ID,
		jwtSecret,
		expiry,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
