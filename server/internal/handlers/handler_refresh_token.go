package handlers

import (
	"database/sql"
	"errors"
	"net/http"
)

func (s *Server) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	// Get decoded refresh token from client
	var refreshToken RefreshToken
	if err := DecodeJson(r, &refreshToken); err != nil {
		RespondWithError(w, s.Logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Get userdata from refresh token
	user, err := s.Store.Queries.GetUserFromRefreshToken(r.Context(), refreshToken.RefreshToken)
	if err != nil {
		// Check if this is a no rows error or actual db error
		if errors.Is(err, sql.ErrNoRows) {
			// Token doesn't exist, is revoked, or is expired
			RespondWithError(w, s.Logger, "Invalid or expired refresh token", err, http.StatusUnauthorized)
			return
		}
		// Actual db error
		RespondWithError(w, s.Logger, "Database error", err, http.StatusInternalServerError)
		return
	}

	// Get fresh access token
	newJwtToken, err := GetJWTToken(user, s.Cfg.JWTSecret, s.Cfg.AccessTokenExpiry)
	if err != nil {
		RespondWithError(w, s.Logger, "Error getting new access token", err, http.StatusInternalServerError)
		return
	}

	// Send new access token as response
	res := AccessTokenResponse{
		AccessToken: newJwtToken,
	}
	if err := RespondWithJSON(w, http.StatusOK, res); err != nil {
		s.Logger.Println("failed to send new access token response:", err)
		return
	}
}
