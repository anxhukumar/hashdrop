package storageguard

import (
	"context"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

// User can only create these many files
const userFilesCountLimit = 500

// Get total sum of bytes consumed by uploaded files and validate if its within limits
func ValidateGlobalS3BucketStorageQuota(ctx context.Context, queries *database.Queries, globalLimit int64) (bool, error) {
	totalBytes, err := queries.GetTotalBytesOfUploadedFiles(ctx)
	if err != nil {
		return false, err
	}

	if totalBytes > globalLimit {
		return false, nil
	}

	return true, nil
}

// Get total sum of bytes consumed by uploaded files of a particular user and validate if its within limits
func ValidateUserS3BucketStorageQuota(ctx context.Context, queries *database.Queries, userID uuid.UUID, userLimit int64) (bool, error) {
	totalBytes, err := queries.GetUsersTotalBytesOfUploadedFiles(ctx, userID)
	if err != nil {
		return false, err
	}

	if totalBytes > userLimit {
		return false, nil
	}

	// check if user has a limited number of files
	numberOfFiles, err := queries.CountFilesOfUser(ctx, userID)
	if err != nil {
		return false, err
	}

	if numberOfFiles > userFilesCountLimit {
		return false, nil
	}

	return true, nil
}
