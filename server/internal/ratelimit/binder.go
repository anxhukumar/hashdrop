package ratelimit

import (
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/handlers"
)

// Binder connects  specific API routes to their respective rate limiters.
// It uses a "Decorator" pattern to wrap handlers in one or more layers of limiting.
type Binder struct {
	Server   *handlers.Server
	Limiters *Limiters
}

// Reset applies a simple global limit to the admin reset endpoint.
func (b *Binder) Reset(next http.Handler) http.Handler {
	return GlobalRateLimit(b.Limiters.ResetGlobalLimiter, b.Server)(next)
}

// Healthz applies a global limit to the health check to prevent DDoS on monitoring.
func (b *Binder) Healthz(next http.Handler) http.Handler {
	return GlobalRateLimit(b.Limiters.HealthzGlobalLimiter, b.Server)(next)
}

// Order: IP Limit checked FIRST, then Global Limit checked SECOND.

// Auth handles Registration and Login.
func (b *Binder) Auth(next http.Handler) http.Handler {
	h := GlobalRateLimit(b.Limiters.AuthGlobalLimiter, b.Server)(next)
	h = IPRateLimit(h, b.Limiters.AuthIPLimiter, b.Server)
	return h
}

// Token handles token refreshing and revocation.
func (b *Binder) Token(next http.Handler) http.Handler {
	h := GlobalRateLimit(b.Limiters.TokenGlobalLimiter, b.Server)(next)
	h = IPRateLimit(h, b.Limiters.TokenIPLimiter, b.Server)
	return h
}

// Upload handles expensive S3/DB operations.
func (b *Binder) Upload(next http.Handler) http.Handler {
	h := GlobalRateLimit(b.Limiters.UploadGlobalLimiter, b.Server)(next)
	h = UserIDRateLimit(h, b.Limiters.UploadUserLimiter, b.Server)
	return h
}

// List handles bulk metadata retrieval.
func (b *Binder) List(next http.Handler) http.Handler {
	h := GlobalRateLimit(b.Limiters.ListGlobalLimiter, b.Server)(next)
	h = UserIDRateLimit(h, b.Limiters.ListUserLimiter, b.Server)
	return h
}

// FileMeta handles individual file actions (Delete, Salt, Hash).
func (b *Binder) FileMeta(next http.Handler) http.Handler {
	h := GlobalRateLimit(b.Limiters.FileMetaGlobalLimiter, b.Server)(next)
	h = UserIDRateLimit(h, b.Limiters.FileMetaUserLimiter, b.Server)
	return h
}
