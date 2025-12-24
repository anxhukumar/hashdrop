package handlers

import (
	"context"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
)

type authUserKey struct{}

func (s *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get access token from header
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, s.logger, "Missing or invalid access token", err, http.StatusUnauthorized)
			return
		}

		// Validate access token
		userID, err := auth.ValidateJWT(token, s.cfg.JWTSecret)
		if err != nil {
			RespondWithError(w, s.logger, "Invalid or expired access token", err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authUserKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
