package auth

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func Register(email, password string) error {

	reqBody := NewUserOutgoing{
		Email:    email,
		Password: password,
	}

	// Struct to receive decoded json response
	respBody := NewUserIncoming{}

	// Post data
	err := api.PostJSON(config.RegisterEndpoint, reqBody, &respBody)
	if err != nil {
		return fmt.Errorf("register failed: %w", err)
	}

	if respBody.Email == "" {
		return errors.New("register failed: invalid response")
	}

	return nil
}
