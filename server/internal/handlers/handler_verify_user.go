package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/otp"
)

func (s *Server) HandlerVerifyUser(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_verify_user")

	// Get decoded incoming user verification data
	var userVerificationData VerifyRequest
	if err := DecodeJson(r, &userVerificationData); err != nil {
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

	// Get unverified user by email
	userData, err := s.Store.Queries.GetUnverifiedUserByEmail(r.Context(), userVerificationData.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msgToDev := "account does not exist"
			msgToClient := "account doesn't exist"
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
		msgToDev := "error fetching unverified user from database"
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
	logger = logger.With("user_id", userData.ID.String())

	// Get OTP hash of user and compare
	storedOtpHash, err := s.Store.Queries.GetOtpHash(r.Context(), userData.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msgToDev := "otp not available or expired"
			msgToClient := "otp not available or expired"
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
		msgToDev := "error fetching otp hash from database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Validate otp
	if !otp.VerifyOTP(userVerificationData.OTP, storedOtpHash, s.Cfg.OtpHashingSecret) {
		msgToDev := "invalid otp provided for user verification"
		msgToClient := "invalid otp"
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

	// Mark user verified if an email match is found
	err = s.Store.Queries.MarkUserVerifiedByEmail(r.Context(), userVerificationData.Email)
	if err != nil {
		msgToDev := "error marking user as verified in database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Delete otp from db once its verified
	err = s.Store.Queries.DeleteOtpByUserID(r.Context(), userData.ID)
	if err != nil {
		logger.Error("failed to delete otp after verification", "err", err)
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("user verified successfully")
}
