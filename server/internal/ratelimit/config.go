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

	// Thins like Register, Login, Account deletion.
	AuthGlobalLimiter *rate.Limiter
	AuthIPLimiter     *keyRateLimiter

	// Refresh token
	TokenGlobalLimiter *rate.Limiter
	TokenIPLimiter     *keyRateLimiter

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
		// 1 request every 30 seconds.
		ResetGlobalLimiter: rate.NewLimiter(rate.Every(30*time.Second), 1),

		// Healthz: Standard uptime monitoring.
		// 5 requests per second global.
		HealthzGlobalLimiter: rate.NewLimiter(rate.Limit(5), 10),

		// Auth (Register/Login/DeleteAccount):
		// Global: 10 per second.
		// IP: 1 request every 5 seconds.
		AuthGlobalLimiter: rate.NewLimiter(rate.Limit(5), 10),
		AuthIPLimiter:     NewKeyRateLimiter(ctx, rate.Limit(0.1), 2),

		// Token (Refresh/Revoke):
		// Frequent but lightweight.
		TokenGlobalLimiter: rate.NewLimiter(rate.Limit(20), 40),
		TokenIPLimiter:     NewKeyRateLimiter(ctx, rate.Limit(2), 10),

		// PRIVATE (S3 / DB INTENSIVE)

		// Upload (Presign/Complete):
		// Global: 10 uploads/sec.
		// User: 1 upload every 2 seconds.
		UploadGlobalLimiter: rate.NewLimiter(rate.Limit(5), 10),
		UploadUserLimiter:   NewKeyRateLimiter(ctx, rate.Limit(0.2), 3),

		// List (GetAllFiles, ResolveMatches):
		// Global: 50/sec (SQLite reads are very fast).
		// User: 5/sec.
		ListGlobalLimiter: rate.NewLimiter(rate.Limit(40), 80),
		ListUserLimiter:   NewKeyRateLimiter(ctx, rate.Limit(5), 10),

		// FileMeta (Detail, Salt, Hash, Delete):
		// Global: 30/sec.
		// User: 10/sec.
		FileMetaGlobalLimiter: rate.NewLimiter(rate.Limit(20), 40),
		FileMetaUserLimiter:   NewKeyRateLimiter(ctx, rate.Limit(2), 5),
	}
}
