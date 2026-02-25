package cleaners

import (
	"context"
	"fmt"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/handlers"
)

// Deletes file metadata from database that are marked 'failed' or 'deleted' after a certain duration
func deleteStaleFileMetadata(ctx context.Context, s *handlers.Server, olderThan time.Duration) {

	cutOffTime := time.Now().UTC().Add(-olderThan)

	// Delete
	err := s.Store.Queries.CleanDeletedAndFailedFiles(ctx, cutOffTime)
	if err != nil {
		err := fmt.Errorf("error deleting failed or deleted marked files metadata from database: %w", err)
		logErrorWhileCleaning(s.Logger, "delete_stale_file_metadata", err)
		return
	}
}

// Deletes revoked and expired refresh tokens older than a certain duration
func deleteStaleRefreshToken(ctx context.Context, s *handlers.Server, olderThan time.Duration) {

	cutOffTime := time.Now().UTC().Add(-olderThan)

	// Delete
	err := s.Store.Queries.CleanRevokedAndExpiredToken(ctx, cutOffTime)
	if err != nil {
		err := fmt.Errorf("error deleting revoked or expired refresh token: %w", err)
		logErrorWhileCleaning(s.Logger, "delete_stale_refresh_token", err)
		return
	}
}

// Deletes download attempts count older than a certain duration
func deleteStaleDownloadAttemptsCount(ctx context.Context, s *handlers.Server, olderThan time.Duration) {

	cutoffDate := time.Now().UTC().Add(-olderThan)

	// Delete
	err := s.Store.Queries.CleanDownloadCount(ctx, cutoffDate)
	if err != nil {
		err := fmt.Errorf("error deleting old download attmepts count: %w", err)
		logErrorWhileCleaning(s.Logger, "delete_stale_download_attempts_count", err)
		return
	}
}

// Deletes unverified user older than a certain duration to avoid deletion before verification window
func deleteStaleUnverifiedUser(ctx context.Context, s *handlers.Server, olderThan time.Duration) {

	cutOffTime := time.Now().UTC().Add(-olderThan)

	// Delete
	err := s.Store.Queries.CleanUnverifiedUser(ctx, cutOffTime)
	if err != nil {
		err := fmt.Errorf("error deleting unverified users: %w", err)
		logErrorWhileCleaning(s.Logger, "delete_stale_unverified_user", err)
		return
	}
}

// Deletes expired otp's
func deleteStaleOtp(ctx context.Context, s *handlers.Server, olderThan time.Duration) {

	cutOffTime := time.Now().UTC().Add(-olderThan)

	// Delete
	err := s.Store.Queries.CleanExpiredOtp(ctx, cutOffTime)
	if err != nil {
		err := fmt.Errorf("error deleting expired otp's: %w", err)
		logErrorWhileCleaning(s.Logger, "delete_stale_otp", err)
		return
	}
}
