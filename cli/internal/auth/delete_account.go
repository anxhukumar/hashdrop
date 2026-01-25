package auth

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
)

func DeleteAccount(email, password string) error {

	reqBody := UserLoginOutgoing{
		Email:    email,
		Password: password,
	}

	// Make delete user request
	err := api.DeleteAccount(&reqBody)
	if err != nil {
		return fmt.Errorf("delete account: %w", err)
	}

	return nil
}
