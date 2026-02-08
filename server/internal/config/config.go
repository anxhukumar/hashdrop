package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                      string
	DbURL                     string
	JWTSecret                 string
	AccessTokenExpiry         time.Duration
	RefreshTokenExpiry        time.Duration
	Platform                  string
	S3BucketRegion            string
	S3PresignedLinkExpiry     time.Duration
	S3PerFileMaxDataSize      int64
	S3Bucket                  string
	UserIDHashSalt            string
	S3GlobalQuotaLimit        int64
	S3UserSpecificQuotaLimit  int64
	CloudfrontURLPrefix       string
	CloudfrontKeyPairID       string
	CloudfrontPrivateKeyPath  string
	DailyPerFileDownloadLimit int
	OtpHashingSecret          string
}

// Load environment variables and return a config struct
func LoadConfig() (*Config, error) {

	godotenv.Load("secrets/.env")

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

	// Convert string to integer
	dailyPerFileDownloadLimitInt, err := strconv.Atoi(getEnv("DAILY_PER_FILE_DOWNLOAD_LIMIT"))
	if err != nil {
		return nil, fmt.Errorf("error while converting DAILY_PER_FILE_DOWNLOAD_LIMIT string to int")
	}

	cfg := &Config{
		Port:                      getEnv("PORT"),
		DbURL:                     getEnv("DB"),
		JWTSecret:                 getEnv("JWT_SECRET"),
		AccessTokenExpiry:         accessTokenExpiry,
		RefreshTokenExpiry:        refreshTokenExpiry,
		Platform:                  getEnv("PLATFORM"),
		S3BucketRegion:            getEnv("S3_BUCKET_REGION"),
		S3PresignedLinkExpiry:     s3PresignedLinkExpiry,
		S3PerFileMaxDataSize:      int64(52_428_800), // 50 MB maximum
		S3Bucket:                  getEnv("S3_BUCKET"),
		UserIDHashSalt:            getEnv("USERID_HASHING_SALT"),
		S3GlobalQuotaLimit:        int64(20_000_000_000), // 20 GB maximum
		S3UserSpecificQuotaLimit:  int64(1_000_000_000),  // 1 GB maximum
		CloudfrontURLPrefix:       getEnv("CLOUDFRONT_URL_PREFIX"),
		CloudfrontKeyPairID:       getEnv("CLOUDFRONT_KEY_PAIR_ID"),
		CloudfrontPrivateKeyPath:  getEnv("CLOUDFRONT_PRIVATE_KEY_PATH"),
		DailyPerFileDownloadLimit: dailyPerFileDownloadLimitInt,
		OtpHashingSecret:          getEnv("OTP_HASHING_SECRET"),
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
		"PORT":                          c.Port,
		"DB":                            c.DbURL,
		"JWT_SECRET":                    c.JWTSecret,
		"PLATFORM":                      c.Platform,
		"S3_BUCKET_REGION":              c.S3BucketRegion,
		"S3_BUCKET":                     c.S3Bucket,
		"USERID_HASHING_SALT":           c.UserIDHashSalt,
		"CLOUDFRONT_URL_PREFIX":         c.CloudfrontURLPrefix,
		"CLOUDFRONT_KEY_PAIR_ID":        c.CloudfrontKeyPairID,
		"CLOUDFRONT_PRIVATE_KEY_PATH":   c.CloudfrontPrivateKeyPath,
		"DAILY_PER_FILE_DOWNLOAD_LIMIT": getEnv("DAILY_PER_FILE_DOWNLOAD_LIMIT"),
		"OTP_HASHING_SECRET":            c.OtpHashingSecret,
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
	if c.S3PerFileMaxDataSize <= 0 {
		return fmt.Errorf("S3_MAX_DATA_SIZE must be positive")
	}
	if c.S3GlobalQuotaLimit <= 0 {
		return fmt.Errorf("S3_GLOBAL_QUOTA_LIMIT must be positive")
	}
	if c.S3UserSpecificQuotaLimit <= 0 {
		return fmt.Errorf("S3_USER_SPECIFIC_QUOTA_LIMIT must be positive")
	}

	return nil
}

// Fetch environment variables
func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}
