package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/aws"
	"github.com/anxhukumar/hashdrop/server/internal/database"
	storageguard "github.com/anxhukumar/hashdrop/server/internal/storage_guard"
	"github.com/google/uuid"
)

func (s *Server) HandlerGeneratePresignLink(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming file metadata
	var FileMetadata FileUploadRequest
	if err := DecodeJson(r, &FileMetadata); err != nil {
		RespondWithError(w, s.Logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, s.Logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	// Check if the total space consumed by uploaded files is within limits
	// Global quota
	valid, err := storageguard.ValidateGlobalS3BucketStorageQuota(r.Context(), s.Store.Queries, s.Cfg.S3GlobalQuotaLimit)
	if err != nil {
		RespondWithError(
			w,
			s.Logger,
			"Internal server error",
			errors.New("error validating global s3 bucket quota"),
			http.StatusInternalServerError,
		)
		return
	}
	if !valid {
		RespondWithError(
			w,
			s.Logger,
			"Uploads are temporarily disabled due to system storage limits",
			errors.New("global storage limit exceeded"),
			http.StatusServiceUnavailable,
		)
		return
	}
	// User specific quota
	valid, err = storageguard.ValidateUserS3BucketStorageQuota(r.Context(), s.Store.Queries, userID, s.Cfg.S3UserSpecificQuotaLimit)
	if err != nil {
		RespondWithError(
			w,
			s.Logger,
			"Internal server error",
			errors.New("error validating users s3 bucket quota"),
			http.StatusInternalServerError,
		)
		return
	}
	if !valid {
		RespondWithError(
			w,
			s.Logger,
			"Storage limit reached",
			errors.New("user storage limit exceeded"),
			http.StatusForbidden,
		)
		return
	}

	// Generate userID hash to use it as prefix in s3
	s3KeyPrefix := GenerateUserIDHash(userID.String(), s.Cfg.UserIDHashSalt)

	// Generate fileID
	fileID := uuid.New()

	// Make object key ("user-<user_id>/<file_id>")
	s3ObjectKey := fmt.Sprintf("usrh-%s/%s", s3KeyPrefix, fileID.String())

	// Generate presigned link with aws s3
	presignedLinkResponse, err := aws.GeneratePresignedPUT(
		r.Context(),
		s.S3Client,
		s.Cfg.S3PresignedLinkExpiry,
		s.Cfg.S3Bucket,
		s3ObjectKey,
	)
	if err != nil {
		RespondWithError(w, s.Logger, "Error generating presigned put link", err, http.StatusInternalServerError)
		return
	}

	// Upload pending file metadata to database
	fileData := database.CreatePendingFileParams{
		ID:       fileID,
		UserID:   userID,
		FileName: FileMetadata.FileName,
		MimeType: sql.NullString{String: FileMetadata.MimeType, Valid: true},
		S3Key:    s3ObjectKey,
	}
	err = s.Store.Queries.CreatePendingFile(r.Context(), fileData)
	if err != nil {
		RespondWithError(w, s.Logger, "Error creating file meta data", err, http.StatusInternalServerError)
		return
	}

	// Send presign data to the client
	resp := PresignResponse{
		FileID:         fileID,
		UploadResource: *presignedLinkResponse,
	}

	if err := RespondWithJSON(w, http.StatusOK, resp); err != nil {
		s.Logger.Println("failed to send response:", err)
		return
	}
}
