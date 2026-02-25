package cleaners

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

const cleanerNameForLogging = "s3_cleaner"

// Deletes pending marked s3 files older than a certain duration
func deletePendingS3File(ctx context.Context, s *handlers.Server, olderThan time.Duration) {

	cutOffTime := time.Now().UTC().Add(-olderThan)

	// Get all pending files older than cutOffTime
	pendingFiles, err := s.Store.Queries.GetStalePendingFiles(ctx, cutOffTime)
	if err != nil {
		err := fmt.Errorf("error while getting pending files from database: %w", err)
		logErrorWhileCleaning(s.Logger, cleanerNameForLogging, err)
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
			var notFound *types.NotFound
			if errors.As(err, &notFound) {
				err := fmt.Errorf("s3 object not found: %w", err)
				logErrorWhileCleaning(s.Logger, cleanerNameForLogging, err)
				// Mark metadata in database as failed
				err = s.Store.Queries.UpdateFailedFile(
					ctx,
					database.UpdateFailedFileParams{
						ID:     file.ID,
						UserID: file.UserID,
					},
				)
				if err != nil {
					err := fmt.Errorf("error while deleting file metadata: %w", err)
					logErrorWhileCleaning(s.Logger, cleanerNameForLogging, err)
					return
				}

				continue

			}

			err := fmt.Errorf("error while getting s3 head object: %w", err)
			logErrorWhileCleaning(s.Logger, cleanerNameForLogging, err)
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
			logErrorWhileCleaning(s.Logger, cleanerNameForLogging, err)
			continue
		}

		// Mark metadata in database as failed
		err = s.Store.Queries.UpdateFailedFile(
			ctx,
			database.UpdateFailedFileParams{
				ID:     file.ID,
				UserID: file.UserID,
			},
		)
		if err != nil {
			err := fmt.Errorf("error while deleting file metadata: %w", err)
			logErrorWhileCleaning(s.Logger, cleanerNameForLogging, err)
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
		"cleaner_name", cleanerNameForLogging,
		"user_id", userID.String(),
		"verified_file_size", verifiedFileSizeString,
	)

	log.Info("file deleted from s3")
}

// Log error while running any cleaner
func logErrorWhileCleaning(logger *slog.Logger, cleanerName string, err error) {

	// Add context about user and error
	logger.Error(
		"error while running cleaners",
		"cleaner_name", cleanerName,
		"err", err,
	)
}
