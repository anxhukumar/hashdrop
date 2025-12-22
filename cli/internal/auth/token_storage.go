package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDirName  = ".hashdrop"   // name of the tokens directory
	TokensFileName = "tokens.json" // name of the tokens file
)

// Get tokens path
func tokensPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ConfigDirName, TokensFileName), nil
}

// Stores the tokens in the users home directory ~/.hashdrop/tokens.json
func StoreTokens(tokens UserLoginIncoming) error {

	// Get path
	tokensPath, err := tokensPath()
	if err != nil {
		return err
	}

	// Create directory
	configDir := filepath.Dir(tokensPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal the data with proper indentation
	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tokens: %w", err)
	}

	// Create the tokens.json
	if err := os.WriteFile(tokensPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// Load tokens
func LoadTokens() (UserLoginIncoming, error) {
	var tokens UserLoginIncoming

	path, err := tokensPath()
	if err != nil {
		return tokens, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return tokens, fmt.Errorf("failed to read tokens file: %w", err)
	}

	if err := json.Unmarshal(data, &tokens); err != nil {
		return tokens, fmt.Errorf("failed to decode tokens json: %w", err)
	}

	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		return tokens, errors.New("invalid token file")
	}

	return tokens, nil
}

// Delete tokens
func DeleteTokens() error {
	path, err := tokensPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove tokens file: %w", err)
	}
	return nil
}
