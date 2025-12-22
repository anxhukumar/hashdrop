package auth

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func RevokeRefreshToken() error {

	// Get tokens
	tokenData, err := LoadTokens()
	if err != nil {
		// Already logged-out locally
		return nil
	}

	reqBody := RefreshToken{
		RefreshToken: tokenData.RefreshToken,
	}

	// Post data
	err = api.PostJSON(config.RevokeRefreshTokenEndpoint, reqBody, nil)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
