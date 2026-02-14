package cloudfrontguard

import (
	"context"
	"testing"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

func TestValidateDownloadAttempts(t *testing.T) {
	ctx := context.Background()

	db := setupTestDB(t)
	queries := database.New(db)

	fileID := uuid.New()

	tests := []struct {
		name     string
		maxLimit int
		wantOK   bool
		wantErr  bool
	}{
		{
			name:     "Under limit",
			maxLimit: 5,
			wantOK:   true,
			wantErr:  false,
		},
		{
			name:     "Over limit",
			maxLimit: 0,
			wantOK:   false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateDownloadAttempts(ctx, queries, tt.maxLimit, fileID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateDownloadAttempts() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if ok != tt.wantOK {
				t.Fatalf("ValidateDownloadAttempts() = %v, want %v", ok, tt.wantOK)
			}
		})
	}
}
