package handlers

import "net/http"

func (s *Server) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	// Get decoded refresh token from client
	var refreshToken RefreshToken
	if err := DecodeJson(r, &refreshToken); err != nil {
		RespondWithError(w, s.logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Set the revoked_at value in refresh_tokens in database
	err := s.store.Queries.RevokeRefreshToken(r.Context(), refreshToken.RefreshToken)
	if err != nil {
		RespondWithError(w, s.logger, "Error revoking refresh token", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
