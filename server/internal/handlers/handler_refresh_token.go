package handlers

import (
	"database/sql"
	"errors"
	"net/http"
)

func (s *Server) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_refresh_token")

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

	// Get userdata from refresh token
	user, err := s.Store.Queries.GetUserFromRefreshToken(r.Context(), refreshToken.RefreshToken)
	if err != nil {
		// Check if this is a no rows error or actual db error
		if errors.Is(err, sql.ErrNoRows) {
			msgToDev := "invalid, expired, or revoked refresh token"
			msgToClient := "invalid or expired refresh token"
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
		// Actual db error
		msgToDev := "database error while fetching user from refresh token"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Attach user_id in logger context
	logger = logger.With("user_id", user.ID.String())

	// Get fresh access token
	newJwtToken, err := GetJWTToken(user, s.Cfg.JWTSecret, s.Cfg.AccessTokenExpiry)
	if err != nil {
		msgToDev := "error getting new access token"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Send new access token as response
	res := AccessTokenResponse{
		AccessToken: newJwtToken,
	}
	if err := RespondWithJSON(w, http.StatusOK, res); err != nil {
		logger.Error("failed to send new access token response", "err", err)
		return
	}
}
