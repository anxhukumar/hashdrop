package storageguard

import (
	"context"
	"database/sql"
	"testing"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

func simulateFileUpload(t *testing.T, size int64, userID uuid.UUID, ctx context.Context, queries *database.Queries) {
	t.Helper()

	// Create pending file
	fileID := uuid.New()

	if userID == uuid.Nil {
		userID = uuid.New()
	}

	if err := queries.CreatePendingFile(
		ctx,
		database.CreatePendingFileParams{
			ID:       fileID,
			UserID:   userID,
			FileName: "test.txt",
			MimeType: sql.NullString{
				String: "text/plain",
				Valid:  true,
			},
			S3Key: "test/key/" + fileID.String(),
		},
	); err != nil {
		t.Fatalf("error creating dummy file upload: %v", err)
	}

	// Mark file as uploaded with given size
	if err := queries.UpdateUploadedFile(
		ctx,
		database.UpdateUploadedFileParams{
			PlaintextHash: sql.NullString{
				String: "dummyhash",
				Valid:  true,
			},
			PlaintextSizeBytes: sql.NullInt64{
				Int64: size,
				Valid: true,
			},
			EncryptedSizeBytes: sql.NullInt64{
				Int64: size,
				Valid: true,
			},
			KeyManagementMode: sql.NullString{
				String: "vault",
				Valid:  true,
			},
			PassphraseSalt: sql.NullString{
				String: "testsalt",
				Valid:  true,
			},
			Status: "uploaded", // or whatever status your schema expects
			ID:     fileID,
			UserID: userID,
		},
	); err != nil {
		t.Fatalf("error updating dummy file upload: %v", err)
	}
}
