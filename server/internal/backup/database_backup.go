package backup

import (
	"context"
	"os"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const processNameForLog = "database_backup"
const fileSizeLimit = 4 * 1024 * 1024 * 1024 // 4GB

func DatabaseBackup(ctx context.Context, s *handlers.Server, interval time.Duration) {
	if interval <= 0 {
		panic("backup interval must be > 0")
	}

	go func() {
		// Run once before running ticker
		dbBackupHandler(ctx, s)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				dbBackupHandler(ctx, s)
			}
		}
	}()

}

func dbBackupHandler(ctx context.Context, s *handlers.Server) {

	// Create temp file of db to avoid collisions if the file is being written while upload
	tmpFile, err := os.CreateTemp("", "hashdrop-backup-*.db")
	if err != nil {
		s.Logger.Error("error while creating temp file while database backup", "err", err, "process", processNameForLog)
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = s.Store.DB.ExecContext(ctx, "VACUUM INTO ?", tmpFile.Name())
	if err != nil {
		s.Logger.Error("failed to create consistent backup", "err", err, "process", processNameForLog)
		return
	}

	// Validate file size of storage.db
	info, err := os.Stat(tmpFile.Name())
	if err != nil {
		s.Logger.Error("error while fetching file info while database backup", "err", err, "process", processNameForLog)
		return
	}
	fileSize := info.Size()

	if fileSize > fileSizeLimit {
		s.Logger.Error("the database file is bigger than the specified file size for backup", "process", processNameForLog)
		return
	}

	// hashdrop-db-backup-2026-03-01-12-31/storage.db
	objectKey := "backup/database-backup-" + time.Now().Format("2006-01-02-15-04") + "/storage.db"

	// REWIND: Move the cursor back to the start so S3 can read the data
	if _, err := tmpFile.Seek(0, 0); err != nil {
		s.Logger.Error("error seeking to start of temp file", "err", err, "process", processNameForLog)
		return
	}

	_, err = s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.Cfg.S3Bucket,
		Key:    &objectKey,
		Body:   tmpFile,
	})
	if err != nil {
		s.Logger.Error("error while uploading database to s3 bucket while taking backup", "err", err, "process", processNameForLog)
		return
	}

	s.Logger.Info("Backed up database successfully")

}
