package ratelimit

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"golang.org/x/time/rate"
)

// KeyCleanupTTL is the maximum age after which inactive visitors entries are removed.
const (
	keyCleanupTTL   = 2 * time.Hour
	cleanupInterval = 10 * time.Minute
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// keys can be: "ip:1.2.3.4" || "uid:12345" || "email:test@example.com"
type keyRateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

func NewKeyRateLimiter(ctx context.Context, r rate.Limit, b int) *keyRateLimiter {
	limiter := &keyRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    b,
	}

	// cleanup old entries
	go func() {

		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				limiter.mu.Lock()
				for key, v := range limiter.visitors {
					if time.Since(v.lastSeen) > keyCleanupTTL {
						delete(limiter.visitors, key)
					}
				}
				limiter.mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	return limiter
}

func (k *keyRateLimiter) getLimiter(key string) *rate.Limiter {
	k.mu.Lock()
	defer k.mu.Unlock()

	v, exist := k.visitors[key]
	if !exist {
		limiter := rate.NewLimiter(k.rate, k.burst)
		k.visitors[key] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	// update timestamp
	v.lastSeen = time.Now()
	return v.limiter
}

// IP Rate Limit Middleware
func IPRateLimit(next http.Handler, k *keyRateLimiter, s *handlers.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string
		var err error

		forwarded := r.Header.Get("X-Forwarded-For")
		if forwarded != "" {
			ip = strings.TrimSpace(strings.Split(forwarded, ",")[0])
		} else {
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				handlers.RespondWithError(w, s.Logger, "invalid IP", errors.New("invalid IP"), http.StatusBadRequest)
				return
			}
		}

		limiter := k.getLimiter("ip:" + ip)
		if !limiter.Allow() {
			handlers.RespondWithError(w, s.Logger, "too many requests", errors.New("ip rate limit"), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UserID Rate Limit Middleware
func UserIDRateLimit(next http.Handler, k *keyRateLimiter, s *handlers.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get userID from context
		userID, ok := handlers.UserIDFromContext(r.Context())
		if !ok {
			handlers.RespondWithError(w, s.Logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
			return
		}

		limiter := k.getLimiter("uid:" + userID.String())
		if !limiter.Allow() {
			handlers.RespondWithError(w, s.Logger, "too many requests", errors.New("user rate limit"), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
