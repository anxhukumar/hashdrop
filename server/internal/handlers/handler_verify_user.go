package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/otp"
)

func (s *Server) HandlerVerifyUser(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming user verification data
	var userVerificationData VerifyRequest
	if err := DecodeJson(r, &userVerificationData); err != nil {
		RespondWithError(w, s.Logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Get unverified user by email
	userData, err := s.Store.Queries.GetUnverifiedUserByEmail(r.Context(), userVerificationData.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(
				w,
				s.Logger,
				"Account doesn't exist",
				errors.New("The account that user is trying to verify doesn't exist"),
				http.StatusUnauthorized)
			return
		}
		RespondWithError(
			w,
			s.Logger,
			"Error while getting user account",
			fmt.Errorf("Error fetching users account while verifying: %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	// Get OTP hash of user and compare
	storedOtpHash, err := s.Store.Queries.GetOtpHash(r.Context(), userData.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(
				w,
				s.Logger,
				"OTP not available or expired",
				errors.New("The otp user is trying to fetch doesn't exist or is expired"),
				http.StatusBadRequest,
			)
			return
		}
		RespondWithError(
			w,
			s.Logger,
			"Error while getting otp hash",
			fmt.Errorf("Error fetching otp hash while verifying user: %w", err),
			http.StatusInternalServerError,
		)
		return

	}

	// Validate otp
	if !otp.VerifyOTP(userVerificationData.OTP, storedOtpHash, s.Cfg.OtpHashingSecret) {
		RespondWithError(w, s.Logger, "Invalid OTP", nil, http.StatusUnauthorized)
		return
	}

	// Mark user verified if an email match is found
	err = s.Store.Queries.MarkUserVerifiedByEmail(r.Context(), userVerificationData.Email)
	if err != nil {
		RespondWithError(
			w,
			s.Logger,
			"Error verifying account",
			fmt.Errorf("Error verifying user account: %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	// Delete otp from db once its verified
	err = s.Store.Queries.DeleteOtpByUserID(r.Context(), userData.ID)
	if err != nil {
		s.Logger.Printf("warning: failed to delete otp: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
