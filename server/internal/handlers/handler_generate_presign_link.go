package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/aws"
	"github.com/anxhukumar/hashdrop/server/internal/database"
	storageguard "github.com/anxhukumar/hashdrop/server/internal/storage_guard"
	"github.com/google/uuid"
)

func (s *Server) HandlerGeneratePresignLink(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_generate_presign_link")

	// Get decoded incoming file metadata
	var FileMetadata FileUploadRequest
	if err := DecodeJson(r, &FileMetadata); err != nil {
		msgToDev := "user posted invalid json data"
		msgToClient := "invalid JSON payload"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			err,
			http.StatusBadRequest,
		)
		return
	}

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		msgToDev := "user id missing in request context"
		RespondWithError(
			w,
			logger,
			msgToDev,
			nil,
			http.StatusInternalServerError,
		)
		return
	}

	// Attach user_id in logger context to enhance logs
	logger = logger.With("user_id", userID.String())

	// Check if the total space consumed by uploaded files is within limits
	// Global quota
	valid, err := storageguard.ValidateGlobalS3BucketStorageQuota(
		r.Context(),
		s.Store.Queries,
		s.Cfg.S3GlobalQuotaLimit,
	)
	if err != nil {
		msgToDev := "error validating global s3 bucket quota"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}
	if !valid {
		msgToDev := "global storage limit exceeded, uploads temporarily disabled"
		RespondWithError(
			w,
			logger,
			msgToDev,
			nil,
			http.StatusServiceUnavailable,
		)
		return
	}

	// User specific quota
	valid, err = storageguard.ValidateUserS3BucketStorageQuota(
		r.Context(),
		s.Store.Queries,
		userID,
		s.Cfg.S3UserSpecificQuotaLimit,
	)
	if err != nil {
		msgToDev := "error validating user s3 bucket quota"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}
	if !valid {
		msgToDev := "user storage limit exceeded"
		msgToClient := "storage limit reached"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			nil,
			http.StatusForbidden,
		)
		return
	}

	// Generate userID hash to use it as prefix in s3
	s3KeyPrefix := GenerateUserIDHash(userID.String(), s.Cfg.UserIDHashSalt)

	// Generate fileID
	fileID := uuid.New()

	// Attach file_id in logger context
	logger = logger.With("file_id", fileID.String())

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
		msgToDev := "error generating presigned put link"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
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
		msgToDev := "error creating pending file metadata in database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Send presign data to the client
	resp := PresignResponse{
		FileID:         fileID,
		UploadResource: *presignedLinkResponse,
	}

	if err := RespondWithJSON(w, http.StatusOK, resp); err != nil {
		logger.Error("failed to send response", "err", err)
		return
	}
}
