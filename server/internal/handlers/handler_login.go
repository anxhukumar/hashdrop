package handlers

import (
	"net/http"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
	"github.com/anxhukumar/hashdrop/server/internal/database"
)

func (s *Server) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_login")

	// Get decoded incoming user login data
	var userLoginIncoming UserLoginIncoming
	if err := DecodeJson(r, &userLoginIncoming); err != nil {
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

	// Check if user is registered and verified to get account details
	userData, err := s.Store.Queries.GetVerifiedUserByEmail(r.Context(), userLoginIncoming.Email)
	if err != nil {
		msgToDev := "invalid username or password during login"
		msgToClient := "invalid username or password"
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

	// Attach user_id in logger context
	logger = logger.With("user_id", userData.ID.String())

	// Check if password is correct
	isMatch, err := auth.CheckPasswordHash(userLoginIncoming.Password, userData.HashedPassword)
	if err != nil {
		msgToDev := "error while verifying password hash during login"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	if !isMatch {
		msgToDev := "password mismatch during login"
		msgToClient := "invalid username or password"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			nil,
			http.StatusUnauthorized,
		)
		return
	}

	// Get JWT token
	jwtToken, err := GetJWTToken(userData, s.Cfg.JWTSecret, s.Cfg.AccessTokenExpiry)
	if err != nil {
		msgToDev := "error creating jwt access token"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Get refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		msgToDev := "error generating refresh token"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
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
		msgToDev := "error creating refresh token in database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Return output to the client
	loginResponse := UserLoginOutgoing{
		AccessToken:  jwtToken,
		RefreshToken: refreshTokenData.Token,
	}

	if err := RespondWithJSON(w, http.StatusOK, loginResponse); err != nil {
		logger.Error("failed to send login response", "err", err)
		return
	}

	logger.Info("user logged in successfully")
}
