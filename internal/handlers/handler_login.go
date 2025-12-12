package handlers

import (
	"net/http"

	"github.com/anxhukumar/hashdrop/internal/auth"
)

func (s *Server) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming user login data
	var userLoginIncoming UserLoginIncoming
	if err := DecodeJson(r, &userLoginIncoming); err != nil {
		RespondWithError(w, s.logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Check if user is registered and get account details
	userData, err := s.store.Queries.GetUserByEmail(r.Context(), userLoginIncoming.Email)
	if err != nil {
		RespondWithError(w, s.logger, "Invalid username or password", err, http.StatusUnauthorized)
		return
	}

	// Check if password is correct
	isMatch, err := auth.CheckPasswordHash(userLoginIncoming.Password, userData.HashedPassword)
	if err != nil {
		RespondWithError(w, s.logger, "Error while verifying password", err, http.StatusInternalServerError)
		return
	}

	if !isMatch {
		RespondWithError(w, s.logger, "Invalid username or password", err, http.StatusUnauthorized)
		return
	}

}
