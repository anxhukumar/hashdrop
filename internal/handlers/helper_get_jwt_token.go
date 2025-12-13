package handlers

import (
	"time"

	"github.com/anxhukumar/hashdrop/internal/auth"
	"github.com/anxhukumar/hashdrop/internal/database"
)

func GetJWTToken(userData database.User, jwtSecret string, expiry time.Duration) (string, error) {

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
