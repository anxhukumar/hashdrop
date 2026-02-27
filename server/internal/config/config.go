package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const envFilePath = "secrets/.env"

type Config struct {
	Port string

	Platform string

	S3PresignedLinkExpiry    time.Duration
	S3PerFileMaxDataSize     int64
	S3Bucket                 string
	S3GlobalQuotaLimit       int64
	S3UserSpecificQuotaLimit int64

	DbURL string

	JWTSecret string

	CloudfrontKeyPairID      string
	CloudfrontPrivateKeyPath string

	UserIDHashSalt              string
	OtpHashingSecret            string
	RefreshTokenHashingSecretV1 string

	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration

	CloudfrontURLPrefix string

	DailyPerFileDownloadLimit int

	AwsRegion string

	CliVersion string
}

// Load environment variables and return a config struct
func LoadConfig() *Config {

	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalf("Error: could not load .env file: %v", err)
	}

	cfg := &Config{
		Port: getEnv("PORT"),

		Platform: getEnv("PLATFORM"),

		S3PresignedLinkExpiry:    15 * time.Minute,
		S3PerFileMaxDataSize:     int64(52_428_800), // 50 MB maximum
		S3Bucket:                 "hashdrop-files",
		S3GlobalQuotaLimit:       int64(20_000_000_000), // 20 GB maximum
		S3UserSpecificQuotaLimit: int64(1_000_000_000),  // 1 GB maximum

		DbURL: getEnv("DB"),

		JWTSecret: getEnv("JWT_SECRET"),

		CloudfrontKeyPairID:      getEnv("CLOUDFRONT_KEY_PAIR_ID"),
		CloudfrontPrivateKeyPath: getEnv("CLOUDFRONT_PRIVATE_KEY_PATH"),

		UserIDHashSalt:              getEnv("USERID_HASHING_SALT"),
		OtpHashingSecret:            getEnv("OTP_HASHING_SECRET"),
		RefreshTokenHashingSecretV1: getEnv("REFRESH_TOKEN_HASHING_SECRET_VERSION_1"),

		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 720 * time.Hour,

		CloudfrontURLPrefix: "https://cdn.hashdrop.dev/",

		DailyPerFileDownloadLimit: 3,

		AwsRegion: getEnv("AWS_REGION"),

		CliVersion: "1.0.0",
	}

	return cfg
}

func getEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s environment variable cannot be empty", key)
	}

	return value
}
