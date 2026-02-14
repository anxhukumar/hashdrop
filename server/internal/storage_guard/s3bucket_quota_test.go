package storageguard

import (
	"context"
	"testing"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	testutil "github.com/anxhukumar/hashdrop/server/internal/test_util"
	"github.com/google/uuid"
)

func TestValidateGlobalS3BucketStorageQuota(t *testing.T) {
	ctx := context.Background()

	db := testutil.SetupTestDB(t)
	queries := database.New(db)

	// simulate file upload
	userID := uuid.New()
	simulateFileUpload(t, int64(4_000_000_000), userID, ctx, queries)

	// Insert a file totalling 4 GB

	tests := []struct {
		name        string
		globalLimit int64
		wantOK      bool
		wantErr     bool
	}{
		{
			name:        "Under limit",
			globalLimit: int64(5_000_000_000), // 5GB MAXIMUM LIMIT
			wantOK:      true,
			wantErr:     false,
		},
		{
			name:        "Over limit",
			globalLimit: int64(3_000_000_000), // 3GB MAXIMUM LIMIT,
			wantOK:      false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateGlobalS3BucketStorageQuota(ctx, queries, tt.globalLimit)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateGlobalS3BucketStorageQuota() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if ok != tt.wantOK {
				t.Fatalf("ValidateGlobalS3BucketStorageQuota() = %v, want = %v", ok, tt.wantOK)
			}
		})
	}
}

func TestValidateUserS3BucketStorageQuota(t *testing.T) {
	ctx := context.Background()

	db := testutil.SetupTestDB(t)
	queries := database.New(db)

	userA := uuid.New()
	userB := uuid.New()

	// Simulate uploads:
	// userA: 4 GB total
	// userB: 2 GB total
	simulateFileUpload(t, int64(4_000_000_000), userA, ctx, queries)
	simulateFileUpload(t, int64(2_000_000_000), userB, ctx, queries)

	tests := []struct {
		name      string
		userID    uuid.UUID
		userLimit int64
		wantOK    bool
		wantErr   bool
	}{
		{
			name:      "Under limit for userA",
			userID:    userA,
			userLimit: int64(5_000_000_000), // 5 GB limit, userA has 4 GB
			wantOK:    true,
			wantErr:   false,
		},
		{
			name:      "Over limit for userA",
			userID:    userA,
			userLimit: int64(3_000_000_000), // 3 GB limit, userA has 4 GB
			wantOK:    false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateUserS3BucketStorageQuota(ctx, queries, tt.userID, tt.userLimit)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateUserS3BucketStorageQuota() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if ok != tt.wantOK {
				t.Fatalf("ValidateUserS3BucketStorageQuota() = %v, want = %v", ok, tt.wantOK)
			}
		})
	}
}
