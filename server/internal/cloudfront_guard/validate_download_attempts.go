package cloudfrontguard

import (
	"context"
	"fmt"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

// Validates if the request for generating a signed url for a fileID is withing limits for the day.
// Returns true if its within limits and false otherwise.
func ValidateDownloadAttempts(ctx context.Context, queries *database.Queries, maxLimit int, FileID uuid.UUID) (bool, error) {
	attempts, err := queries.CheckAndUpdateDownloadAttemptsCount(
		ctx,
		database.CheckAndUpdateDownloadAttemptsCountParams{
			ID:     uuid.New(),
			FileID: FileID,
		},
	)
	if err != nil {
		return false, fmt.Errorf("error while checking download attempts: %w", err)
	}

	if attempts > int64(maxLimit) {
		return false, nil
	}

	return true, nil
}
