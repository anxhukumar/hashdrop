package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

// Load environment variables and returns a config struct
func LoadConfig() (*Config, error) {

	godotenv.Load()

	cfg := &Config{
		Port: getEnv("PORT"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Check if the configuration is valid
func (c *Config) Validate() error {
	// Invalid port error
	if c.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	return nil
}

// Fetch environment variables
func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}
