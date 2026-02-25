package cleaners

import (
	"context"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/handlers"
)

// Schedule automated cleaners that run in fixed durations constantly
func ScheduledCleaners(
	ctx context.Context,
	s *handlers.Server,
) {
	go cleaner(
		ctx,
		s,
		deletePendingS3File,
		cleanerConfig.pendingFileS3MaxAge,
		cleanerConfig.pendingFileS3Interval,
	)
	go cleaner(
		ctx,
		s,
		deleteStaleFileMetadata,
		cleanerConfig.staleFileMetadataMaxAge,
		cleanerConfig.staleFileMetadataInterval,
	)
	go cleaner(
		ctx,
		s,
		deleteStaleRefreshToken,
		cleanerConfig.staleRefreshTokenMaxAge,
		cleanerConfig.staleRefreshTokenInterval,
	)
	go cleaner(
		ctx,
		s,
		deleteStaleDownloadAttemptsCount,
		cleanerConfig.staleDownloadAttemptsCountMaxAge,
		cleanerConfig.staleDownloadAttemptsCountInterval,
	)
	go cleaner(
		ctx,
		s,
		deleteStaleUnverifiedUser,
		cleanerConfig.staleUnverifiedUserMaxAge,
		cleanerConfig.staleUnverifiedUserInterval,
	)
	go cleaner(
		ctx,
		s,
		deleteStaleOtp,
		cleanerConfig.staleOtpMaxAge,
		cleanerConfig.staleOtpInterval,
	)
}

// Configurations for cleaners
var cleanerConfig = struct {
	pendingFileS3MaxAge                time.Duration
	pendingFileS3Interval              time.Duration
	staleFileMetadataMaxAge            time.Duration
	staleFileMetadataInterval          time.Duration
	staleRefreshTokenMaxAge            time.Duration
	staleRefreshTokenInterval          time.Duration
	staleDownloadAttemptsCountMaxAge   time.Duration
	staleDownloadAttemptsCountInterval time.Duration
	staleUnverifiedUserMaxAge          time.Duration
	staleUnverifiedUserInterval        time.Duration
	staleOtpMaxAge                     time.Duration
	staleOtpInterval                   time.Duration
}{
	pendingFileS3MaxAge:                30 * time.Minute,
	pendingFileS3Interval:              10 * time.Minute,
	staleFileMetadataMaxAge:            10 * time.Minute,
	staleFileMetadataInterval:          10 * time.Minute,
	staleRefreshTokenMaxAge:            10 * time.Minute,
	staleRefreshTokenInterval:          10 * time.Minute,
	staleDownloadAttemptsCountMaxAge:   24 * time.Hour,
	staleDownloadAttemptsCountInterval: 12 * time.Hour,
	staleUnverifiedUserMaxAge:          30 * time.Minute,
	staleUnverifiedUserInterval:        10 * time.Minute,
	staleOtpMaxAge:                     0,
	staleOtpInterval:                   10 * time.Minute,
}

// Calls the cleaner functions on regular intervals.
// It blocks until ctx is cancelled and should be started in a go routine.
func cleaner(
	ctx context.Context,
	s *handlers.Server,
	cleanerFunction func(context.Context, *handlers.Server, time.Duration),
	olderThan time.Duration,
	interval time.Duration,
) {
	if interval <= 0 {
		panic("cleaner interval must be > 0")
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleanerFunction(ctx, s, olderThan)
		}
	}
}
