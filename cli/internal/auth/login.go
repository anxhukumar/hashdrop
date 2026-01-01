package auth

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func Login(email, password string) error {

	reqBody := UserLoginOutgoing{
		Email:    email,
		Password: password,
	}

	// Struct to receive decoded json response
	respBody := UserLoginIncoming{}

	// Post data
	err := api.PostJSON(config.LoginEndpoint, reqBody, &respBody, "")
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	if respBody.AccessToken == "" || respBody.RefreshToken == "" {
		return errors.New("login failed: invalid response")
	}

	// Store token
	if err := StoreTokens(respBody); err != nil {
		return fmt.Errorf("store tokens: %w", err)
	}

	return nil
}
