package cleaners

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Calls the deletePendingS3File() on regular intervals.
// It blocks until ctx is cancelled and should be started in a go routine.
func S3Cleaner(ctx context.Context, s *handlers.Server, cleaningRate time.Duration, olderThan time.Duration) error {
	ticker := time.NewTicker(cleaningRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			deletePendingS3File(ctx, s, olderThan)
		}
	}
}

func deletePendingS3File(ctx context.Context, s *handlers.Server, olderThan time.Duration) {

	cutOffTime := time.Now().UTC().Add(-olderThan)

	// Get all pending files older than cutOffTime
	pendingFiles, err := s.Store.Queries.GetStalePendingFiles(ctx, cutOffTime)
	if err != nil {
		err := fmt.Errorf("error while getting pending files from database: %w", err)
		logErrorFileDeletion(s.Logger, err)
		return
	}

	if len(pendingFiles) == 0 {
		return
	}

	// Delete all pending files
	for _, file := range pendingFiles {

		// Get verified file size
		head, err := s.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(s.Cfg.S3Bucket),
			Key:    aws.String(file.S3Key),
		})
		if err != nil {
			err := fmt.Errorf("error while getting s3 head object: %w", err)
			logErrorFileDeletion(s.Logger, err)
			continue
		}

		verifiedFileSize := aws.ToInt64(head.ContentLength)

		// Delete
		_, err = s.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.Cfg.S3Bucket),
			Key:    aws.String(file.S3Key),
		})
		if err != nil {
			err := fmt.Errorf("error while deleting s3 object: %w", err)
			logErrorFileDeletion(s.Logger, err)
			continue
		}

		// logging
		logPendingFileDeletion(file.UserID, verifiedFileSize, s.Logger)

	}
}

// Log data about deleted file
func logPendingFileDeletion(userID uuid.UUID, verifiedFileSize int64, logger *slog.Logger) {

	// verified file size string
	verifiedFileSizeString := strconv.FormatInt(verifiedFileSize, 10)

	// Add context about file and user
	log := logger.With(
		"cleaner", "s3_cleaner",
		"user_id", userID.String(),
		"verified_file_size", verifiedFileSizeString,
	)

	log.Info("file deleted from s3")
}

// Log error while deleting file
func logErrorFileDeletion(logger *slog.Logger, err error) {

	// Add context about user and error
	logger.Error(
		"error while running s3_cleaner",
		"cleaner", "s3_cleaner",
		"err", err,
	)
}
