package ratelimit

import "golang.org/x/time/rate"

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
