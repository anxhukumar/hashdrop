package ratelimit

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type Limiters struct {
	// ---------------------------------------
	// PUBLIC
	// ---------------------------------------

	// Reset the db data (in dev environment)
	ResetGlobalLimiter *rate.Limiter

	// Server health check, returns OK.
	HealthzGlobalLimiter *rate.Limiter

	// Things like Register, Login, Account deletion.
	AuthGlobalLimiter *rate.Limiter
	AuthIPLimiter     *keyRateLimiter

	// Control otp
	OTPRateGlobalLimiter *rate.Limiter
	OTPRateIPLimiter     *keyRateLimiter

	// Refresh token
	TokenGlobalLimiter *rate.Limiter
	TokenIPLimiter     *keyRateLimiter

	// CLI version check
	CliVersionCheckGlobalLimiter *rate.Limiter
	CliVersionCheckIpLimiter     *keyRateLimiter

	// ---------------------------------------
	// PRIVATE
	// ---------------------------------------

	// Limit file Upload
	// Expensive operation
	UploadGlobalLimiter *rate.Limiter
	UploadUserLimiter   *keyRateLimiter

	// Used with things like listing all files.
	// Expected to be requested a lot.
	ListGlobalLimiter *rate.Limiter
	ListUserLimiter   *keyRateLimiter

	// Used for file specific tasks.
	// Like opening the details of a file, deleting a file, etc.
	FileMetaGlobalLimiter *rate.Limiter
	FileMetaUserLimiter   *keyRateLimiter
}

func NewDefaultLimiters(ctx context.Context) *Limiters {
	return &Limiters{
		// PUBLIC / ADMIN

		// Reset: Extreme protection. Only required for dev.
		// 1 request every 30 seconds. Burst: 1.
		ResetGlobalLimiter: rate.NewLimiter(rate.Every(30*time.Second), 1),

		// Healthz:
		// 5 requests per second globally. Burst: 10.
		HealthzGlobalLimiter: rate.NewLimiter(rate.Limit(5), 10),

		// Auth (Register/Login/DeleteAccount):
		// Global: 5 requests per second. Burst: 10.
		// IP: Refills at 3 requests per minute. Burst: 10.
		// (An IP can send up to 10 immediately, then sustains at 3/minute.)
		AuthGlobalLimiter: rate.NewLimiter(rate.Limit(10), 30),
		AuthIPLimiter:     NewKeyRateLimiter(ctx, rate.Limit(3.0/60.0), 5),

		// OTP (Send/Verify):
		// Global: Refills at 100 requests per hour. Burst: 100.
		// IP: Refills at 3 requests per hour. Burst: 3.
		OTPRateGlobalLimiter: rate.NewLimiter(rate.Limit(100.0/3600.0), 100),
		OTPRateIPLimiter:     NewKeyRateLimiter(ctx, rate.Limit(3.0/3600.0), 3),

		// Token (Refresh/Revoke):
		// Global: 20 requests per second. Burst: 40.
		// IP: 1 request per second. Burst: 5.
		TokenGlobalLimiter: rate.NewLimiter(rate.Limit(10), 30),
		TokenIPLimiter:     NewKeyRateLimiter(ctx, rate.Limit(1), 5),

		// CLI Version Check:
		// Global: 100 requests per second. Burst: 200.
		// IP: 10 requests per second. Burst: 20.
		CliVersionCheckGlobalLimiter: rate.NewLimiter(rate.Limit(50), 100),
		CliVersionCheckIpLimiter:     NewKeyRateLimiter(ctx, rate.Limit(10), 20),

		// PRIVATE (S3 / DB INTENSIVE)

		// Upload (Presign/Complete):
		// Global: 5 uploads per second. Burst: 10.
		// User: Refills at 0.2 requests per second (~1 upload every 5 seconds). Burst: 1.
		UploadGlobalLimiter: rate.NewLimiter(rate.Limit(20), 40),
		UploadUserLimiter:   NewKeyRateLimiter(ctx, rate.Limit(0.2), 1),

		// List (GetAllFiles, ResolveMatches):
		// Global: 40 requests per second. Burst: 80.
		// User: 5 requests per second. Burst: 10.
		ListGlobalLimiter: rate.NewLimiter(rate.Limit(20), 100),
		ListUserLimiter:   NewKeyRateLimiter(ctx, rate.Limit(5), 10),

		// FileMeta (Detail, Salt, Hash, Delete):
		// Global: 20 requests per second. Burst: 40.
		// User: 5 requests per second. Burst: 5.
		FileMetaGlobalLimiter: rate.NewLimiter(rate.Limit(20), 100),
		FileMetaUserLimiter:   NewKeyRateLimiter(ctx, rate.Limit(5), 5),
	}
}
