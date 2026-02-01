package storageguard

import (
	"context"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

// Get total sum of bytes consumed by uploaded files and validate if its within limits
func ValidateGlobalS3BucketStorageQuota(ctx context.Context, queries *database.Queries, globalLimit int64) (bool, error) {
	totalBytes, err := queries.GetTotalBytesOfUploadedFiles(ctx)
	if err != nil {
		return false, err
	}

	// Check if the total bytes are within specified limits
	if totalBytes < globalLimit {
		return true, nil
	}

	return false, nil
}

// Get total sum of bytes consumed by uploaded files of a particular user and validate if its within limits
func ValidateUserS3BucketStorageQuota(ctx context.Context, queries *database.Queries, userID uuid.UUID, userLimit int64) (bool, error) {
	totalBytes, err := queries.GetUsersTotalBytesOfUploadedFiles(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check if the total bytes are within specified limits
	if totalBytes < userLimit {
		return true, nil
	}

	return false, nil
}
