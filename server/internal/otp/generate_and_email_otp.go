package otp

import (
	"context"
	"fmt"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

// This function generates a random otp -> hashes the otp -> saves the hashed otp in database
// -> Emails the otp to users email address.
func GenerateAndEmailOtp(ctx context.Context, userID uuid.UUID, userEmail string, otpHashingSecret string, queries *database.Queries) error {

	originalOtp, err := GenerateOTP()
	if err != nil {
		return fmt.Errorf("error generating random otp: %w", err)
	}

	hashedOtp := HashOTP(originalOtp, otpHashingSecret)

	// Create otp record in database
	otpID := uuid.New()
	err = queries.CreateOtp(ctx, database.CreateOtpParams{
		ID:      otpID,
		UserID:  userID,
		OtpHash: hashedOtp,
	})
	if err != nil {
		return fmt.Errorf("Error creating otp record in database: %w", err)
	}

	// Send otp to users email address
	sender, err := NewSender(ctx)
	if err != nil {
		return fmt.Errorf("Error creating config for aws ses: %w", err)
	}

	err = sender.SendOTP(ctx, userEmail, originalOtp)
	if err != nil {
		// Delete the otp record from database as that will be useless now
		err = queries.DeleteOtpByOtpID(ctx, otpID)
		if err != nil {
			return fmt.Errorf("Error while deleting otp from database using otp id: %w", err)
		}
		return fmt.Errorf("Error while sending otp to users email: %w", err)
	}

	return nil
}
