package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/golang-jwt/jwt"
)

// Safety buffer (in seconds) to avoid using tokens that may expire mid-request.
const tokenExpiryLeeway = 30

// Refreshes an access token (if expired) and retuns access token
func EnsureAccessToken() (string, error) {
	// Get tokens
	tokenData, err := LoadTokens()
	if err != nil {
		return "", err
	}

	// Return the current existing access token if it is valid
	if !isAccessTokenExpired(tokenData.AccessToken) {
		return tokenData.AccessToken, nil
	}

	// Refresh access token
	newAccessToken, err := refreshAccessToken(tokenData.RefreshToken)
	if err != nil {
		// Refresh token invalid or expired
		_ = DeleteTokens()

		return "", errors.New("session expired, please login again")
	}
	newTokens := UserLoginIncoming{
		AccessToken:  newAccessToken,
		RefreshToken: tokenData.RefreshToken,
	}

	// Store the new access token
	if err := StoreTokens(newTokens); err != nil {
		return "", fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return newTokens.AccessToken, nil
}

// Returns true if the access token is expired
func isAccessTokenExpired(token string) bool {
	claims := jwt.MapClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(token, claims)
	if err != nil {
		return true // take error as expired
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return true // take the absence of exp as expired
	}

	return time.Now().Unix() > int64(exp)-tokenExpiryLeeway
}

// Refresh access tokens using refresh tokens
func refreshAccessToken(refreshToken string) (string, error) {

	reqBody := RefreshToken{
		RefreshToken: refreshToken,
	}

	// Struct to receive decoded json response
	respBody := struct {
		AccessToken string `json:"access_token"`
	}{}

	// Post data
	err := api.PostJSON(config.RefreshTokenEndpoint, reqBody, &respBody, "")
	if err != nil {
		return "", fmt.Errorf("refresh access token: %w", err)
	}

	if respBody.AccessToken == "" {
		return "", errors.New("token refresh failed: invalid response")
	}

	return respBody.AccessToken, nil
}
