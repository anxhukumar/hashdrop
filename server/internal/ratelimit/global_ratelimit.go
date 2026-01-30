package ratelimit

import (
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"golang.org/x/time/rate"
)

func GlobalRateLimit(limiter *rate.Limiter, s *handlers.Server) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				handlers.RespondWithError(w, s.Logger, "too many requests", errors.New("too many requests"), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
