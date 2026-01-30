package handlers

import (
	"net/http"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
	"github.com/anxhukumar/hashdrop/server/internal/database"
)

func (s *Server) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming user login data
	var userLoginIncoming UserLoginIncoming
	if err := DecodeJson(r, &userLoginIncoming); err != nil {
		RespondWithError(w, s.Logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Check if user is registered and get account details
	userData, err := s.Store.Queries.GetUserByEmail(r.Context(), userLoginIncoming.Email)
	if err != nil {
		RespondWithError(w, s.Logger, "Invalid username or password", err, http.StatusUnauthorized)
		return
	}

	// Check if password is correct
	isMatch, err := auth.CheckPasswordHash(userLoginIncoming.Password, userData.HashedPassword)
	if err != nil {
		RespondWithError(w, s.Logger, "Error verifying password", err, http.StatusInternalServerError)
		return
	}

	if !isMatch {
		RespondWithError(w, s.Logger, "Invalid username or password", nil, http.StatusUnauthorized)
		return
	}

	// Get JWT token
	jwtToken, err := GetJWTToken(userData, s.Cfg.JWTSecret, s.Cfg.AccessTokenExpiry)
	if err != nil {
		RespondWithError(w, s.Logger, "Error creating auth token", err, http.StatusInternalServerError)
		return
	}

	// Get refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, s.Logger, "Error generating refresh token", err, http.StatusInternalServerError)
		return
	}

	// Send the refresh token to database

	expiry := time.Now().Add(s.Cfg.RefreshTokenExpiry)

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userData.ID,
		ExpiresAt: expiry,
	}
	refreshTokenData, err := s.Store.Queries.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		RespondWithError(w, s.Logger, "Error creating refresh token", err, http.StatusInternalServerError)
		return
	}

	// Return output to the client
	loginResponse := UserLoginOutgoing{
		AccessToken:  jwtToken,
		RefreshToken: refreshTokenData.Token,
	}

	if err := RespondWithJSON(w, http.StatusOK, loginResponse); err != nil {
		s.Logger.Println("failed to send login response:", err)
		return
	}
}
