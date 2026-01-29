package handlers

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ipCleanupTTL is the maximum age after which inactive IP entries are removed.
const (
	keyCleanupTTL   = 10 * time.Minute
	cleanupInterval = time.Minute
)

// GLOBAL RATE LIMIT
func (s *Server) GlobalRateLimit(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				RespondWithError(w, s.logger, "too many requests", errors.New("too many requests"), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// KEY RATE LIMIT
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

	// cleanup old IPs
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
func (s *Server) IPRateLimit(next http.Handler, k *keyRateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: enable X-Forwarded-For when behind a trusted proxy
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			RespondWithError(w, s.logger, "invalid IP", errors.New("invalid IP"), http.StatusBadRequest)
			return
		}

		limiter := k.getLimiter("ip:" + ip)
		if !limiter.Allow() {
			RespondWithError(w, s.logger, "too many requests", errors.New("ip rate limit"), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UserID Rate Limit Middleware
func (s *Server) UserIDRateLimit(next http.Handler, k *keyRateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get userID from context
		userID, ok := UserIDFromContext(r.Context())
		if !ok {
			RespondWithError(w, s.logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
			return
		}

		limiter := k.getLimiter("uid:" + userID.String())
		if !limiter.Allow() {
			RespondWithError(w, s.logger, "too many requests", errors.New("user rate limit"), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
