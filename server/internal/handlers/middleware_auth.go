package handlers

import (
	"context"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
)

type authUserKey struct{}

func (s *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.Logger.With("middleware", "auth")

		// Get access token from header
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			msgToDev := "missing or invalid authorization header"
			msgToClient := "missing or invalid access token"
			RespondWithWarn(
				w,
				logger,
				msgToDev,
				msgToClient,
				err,
				http.StatusUnauthorized,
			)
			return
		}

		// Validate access token
		userID, err := auth.ValidateJWT(token, s.Cfg.JWTSecret)
		if err != nil {
			msgToDev := "invalid or expired jwt access token"
			msgToClient := "invalid or expired access token"
			RespondWithWarn(
				w,
				logger,
				msgToDev,
				msgToClient,
				err,
				http.StatusUnauthorized,
			)
			return
		}

		ctx := context.WithValue(r.Context(), authUserKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
