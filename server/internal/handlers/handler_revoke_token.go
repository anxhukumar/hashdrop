package handlers

import (
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
)

func (s *Server) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_revoke_token")

	// Get decoded refresh token from client
	var refreshToken RefreshToken
	if err := DecodeJson(r, &refreshToken); err != nil {
		msgToDev := "user posted invalid json data"
		msgToClient := "invalid JSON payload"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			err,
			http.StatusBadRequest,
		)
		return
	}

	// Hash refresh token
	hashedRefreshToken := auth.HashRefreshToken(refreshToken.RefreshToken, []byte(s.Cfg.RefreshTokenHashingSecretV1))

	// Set the revoked_at value in refresh_tokens in database
	err := s.Store.Queries.RevokeRefreshToken(r.Context(), hashedRefreshToken)
	if err != nil {
		msgToDev := "error revoking refresh token in database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
