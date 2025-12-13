package handlers

import (
	"net/http"
)

func (s *Server) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	// Get decoded refresh token from client
	var refreshToken RefreshToken
	if err := DecodeJson(r, &refreshToken); err != nil {
		RespondWithError(w, s.logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Get userdata from refresh token
	user, err := s.store.Queries.GetUserFromRefreshToken(r.Context(), refreshToken.RefreshToken)
	if err != nil {
		RespondWithError(w, s.logger, "Invalid or expired refresh token", err, http.StatusUnauthorized)
		return
	}

	// Get fresh access token
	newJwtToken, err := GetJWTToken(user, s.cfg.JWTSecret, s.cfg.AccessTokenExpiry)
	if err != nil {
		RespondWithError(w, s.logger, "Error getting new access token", err, http.StatusInternalServerError)
		return
	}

	// Send new access token as response
	res := AccessTokenResponse{
		AccessToken: newJwtToken,
	}
	if err := RespondWithJSON(w, http.StatusOK, res); err != nil {
		s.logger.Println("failed to send new access token response:", err)
	}
}
