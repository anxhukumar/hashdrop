package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DbURL              string
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Platform           string
}

// Load environment variables and return a config struct
func LoadConfig() (*Config, error) {

	godotenv.Load()

	// Parse durations
	accessTokenExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TOKEN_EXPIRY: %w", err)
	}

	refreshTokenExpiry, err := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRY: %w", err)
	}

	cfg := &Config{
		Port:               getEnv("PORT"),
		DbURL:              getEnv("DB"),
		JWTSecret:          getEnv("JWT_SECRET"),
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
		Platform:           getEnv("PLATFORM"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Check if the configuration is valid
func (c *Config) Validate() error {

	// Maps to each port value for error messages
	checks := map[string]string{
		"PORT":       c.Port,
		"DB":         c.DbURL,
		"JWT_SECRET": c.JWTSecret,
		"PLATFORM":   c.Platform,
	}

	for name, value := range checks {
		if value == "" {
			return fmt.Errorf("%v environment variable cannot be empty", name)
		}
	}

	// Validate durations
	if c.AccessTokenExpiry <= 0 {
		return fmt.Errorf("JWT_ACCESS_TOKEN_EXPIRY must be positive")
	}
	if c.RefreshTokenExpiry <= 0 {
		return fmt.Errorf("REFRESH_TOKEN_EXPIRY must be positive")
	}

	return nil
}

// Fetch environment variables
func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}
