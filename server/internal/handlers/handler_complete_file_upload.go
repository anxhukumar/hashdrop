package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *Server) HandlerCompleteFileUpload(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming file json data
	var FileUploadMetadata FileUploadMetadata
	if err := DecodeJson(r, &FileUploadMetadata); err != nil {
		RespondWithError(w, s.logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, s.logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	// Get metadata of the uploaded file

	// Fetch s3ObjectKey from db
	ObjectKey, err := s.store.Queries.GetS3KeyFromFileID(r.Context(), database.GetS3KeyFromFileIDParams{
		ID:     FileUploadMetadata.FileID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, s.logger, "File not found", err, http.StatusNotFound)
			return
		}
		RespondWithError(w, s.logger, "Error fetching S3ObjectKey from database", err, http.StatusInternalServerError)
		return
	}

	head, err := s.s3Client.HeadObject(r.Context(), &s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.S3Bucket),
		Key:    aws.String(ObjectKey),
	})
	if err != nil {
		RespondWithError(w, s.logger, "Error fetching object metadata from s3", err, http.StatusInternalServerError)
		return
	}

	verifiedFileSize := aws.ToInt64(head.ContentLength) // int64
	if verifiedFileSize > s.cfg.S3MaxDataSize {
		// Delete object
		_, _ = s.s3Client.DeleteObject(r.Context(), &s3.DeleteObjectInput{
			Bucket: aws.String(s.cfg.S3Bucket),
			Key:    aws.String(ObjectKey),
		})

		// Update db status to failed
		if err := s.store.Queries.UpdateFailedFile(r.Context(), database.UpdateFailedFileParams{
			Status: "failed",
			ID:     FileUploadMetadata.FileID,
		}); err != nil {
			RespondWithError(w, s.logger, "Error marking file as failed", err, http.StatusInternalServerError)
			return
		}

		// Respond with error
		RespondWithError(w, s.logger, "File size exceeds the allowed limit", err, http.StatusRequestEntityTooLarge)
		return
	}

	var keyManagementModeVal string
	if FileUploadMetadata.PassphraseSalt == "" {
		keyManagementModeVal = "vault"
	} else {
		keyManagementModeVal = "passphrase"
	}

	// Send data to db
	fileData := database.UpdateUploadedFileParams{
		PlaintextHash: sql.NullString{
			String: FileUploadMetadata.PlaintextHash,
			Valid:  FileUploadMetadata.PlaintextHash != "",
		},
		PlaintextSizeBytes: sql.NullInt64{
			Int64: FileUploadMetadata.PlaintextSizeBytes,
			Valid: FileUploadMetadata.PlaintextSizeBytes > 0,
		},
		EncryptedSizeBytes: sql.NullInt64{
			Int64: verifiedFileSize,
			Valid: verifiedFileSize > 0,
		},
		KeyManagementMode: sql.NullString{
			String: keyManagementModeVal,
			Valid:  keyManagementModeVal != "",
		},
		PassphraseSalt: sql.NullString{
			String: FileUploadMetadata.PassphraseSalt,
			Valid:  FileUploadMetadata.PassphraseSalt != "",
		},
		Status: "uploaded",
		ID:     FileUploadMetadata.FileID,
		UserID: userID,
	}
	if err := s.store.Queries.UpdateUploadedFile(r.Context(), fileData); err != nil {
		RespondWithError(w, s.logger, "Error updating file meta data", err, http.StatusInternalServerError)
		return
	}

	RespondWithJSON(
		w,
		http.StatusOK,
		FileUploadSuccessResponse{
			S3ObjectKey:      ObjectKey,
			UploadedFileSize: verifiedFileSize,
		})

}
