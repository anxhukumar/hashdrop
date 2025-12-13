package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DbURL             string
	JWTSecret         string
	AccessTokenExpiry string
	Platform          string
}

// Load environment variables and return a config struct
func LoadConfig() (*Config, error) {

	godotenv.Load()

	cfg := &Config{
		Port:              getEnv("PORT"),
		DbURL:             getEnv("DB"),
		JWTSecret:         getEnv("JWT_SECRET"),
		AccessTokenExpiry: getEnv("JWT_ACCESS_TOKEN_EXPIRY"),
		Platform:          getEnv("PLATFORM"),
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
		"PORT":                    c.Port,
		"DB":                      c.DbURL,
		"JWT_SECRET":              c.JWTSecret,
		"JWT_ACCESS_TOKEN_EXPIRY": c.AccessTokenExpiry,
		"PLATFORM":                c.Platform,
	}

	for name, value := range checks {
		if value == "" {
			return fmt.Errorf("%v environment variable cannot be empty", name)
		}
	}

	return nil
}

// Fetch environment variables
func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}
