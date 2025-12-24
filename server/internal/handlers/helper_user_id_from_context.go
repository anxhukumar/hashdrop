package handlers

import (
	"context"

	"github.com/google/uuid"
)

// Fetches userID set in the auth middleware
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(authUserKey{}).(uuid.UUID)
	return id, ok
}
