package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Stores the tokens in the users home directory ~/.hashdrop/tokens.json
func StoreTokens(tokens UserLoginIncoming) error {

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create directory
	configDir := filepath.Join(homeDir, ".hashdrop")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal the data with proper indentation
	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tokens: %w", err)
	}

	// Create the tokens.json
	tokensPath := filepath.Join(configDir, "tokens.json")
	if err := os.WriteFile(tokensPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}
