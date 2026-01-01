package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	DbURL                 string
	JWTSecret             string
	AccessTokenExpiry     time.Duration
	RefreshTokenExpiry    time.Duration
	Platform              string
	S3BucketRegion        string
	S3PresignedLinkExpiry time.Duration
	S3MaxDataSize         int64
	S3Bucket              string
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

	s3PresignedLinkExpiry, err := time.ParseDuration(getEnv("S3_PRESIGNED_LINK_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("invalid S3_PRESIGNED_LINK_EXPIRY: %w", err)
	}

	cfg := &Config{
		Port:                  getEnv("PORT"),
		DbURL:                 getEnv("DB"),
		JWTSecret:             getEnv("JWT_SECRET"),
		AccessTokenExpiry:     accessTokenExpiry,
		RefreshTokenExpiry:    refreshTokenExpiry,
		Platform:              getEnv("PLATFORM"),
		S3BucketRegion:        getEnv("S3_BUCKET_REGION"),
		S3PresignedLinkExpiry: s3PresignedLinkExpiry,
		S3MaxDataSize:         int64(52_428_800), // 50 MB maximum
		S3Bucket:              getEnv("S3_BUCKET"),
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
		"PORT":             c.Port,
		"DB":               c.DbURL,
		"JWT_SECRET":       c.JWTSecret,
		"PLATFORM":         c.Platform,
		"S3_BUCKET_REGION": c.S3BucketRegion,
		"S3_BUCKET":        c.S3Bucket,
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
	if c.S3PresignedLinkExpiry <= 0 {
		return fmt.Errorf("S3_PRESIGNED_LINK_EXPIRY must be positive")
	}

	// Validate S3 byte limits
	if c.S3MaxDataSize <= 0 {
		return fmt.Errorf("S3_MAX_DATA_SIZE must be positive")
	}

	return nil
}

// Fetch environment variables
func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}
