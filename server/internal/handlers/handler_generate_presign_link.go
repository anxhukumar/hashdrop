package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/aws"
	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

func (s *Server) HandlerGeneratePresignLink(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming file metadata
	var FileMetadata FileUploadRequest
	if err := DecodeJson(r, &FileMetadata); err != nil {
		RespondWithError(w, s.logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, s.logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	// Generate fileID
	fileID := uuid.New()

	// Make object key ("user-<user_id>/<file_id>")
	s3ObjectKey := fmt.Sprintf("user-%s/%s", userID.String(), fileID.String())

	// Generate presigned link with aws s3
	presignedLinkResponse, err := aws.GeneratePresignedPOST(
		r.Context(),
		s.s3Config,
		s.cfg.S3MinDataSize,
		s.cfg.S3MaxDataSize,
		s.cfg.S3PresignedLinkExpiry,
		s.cfg.S3Bucket,
		s3ObjectKey,
		s.cfg.S3BucketRegion,
	)
	if err != nil {
		RespondWithError(w, s.logger, "Error generating presigned post link", err, http.StatusInternalServerError)
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
	_, err = s.store.Queries.CreatePendingFile(r.Context(), fileData)
	if err != nil {
		RespondWithError(w, s.logger, "Error creating file meta data", err, http.StatusInternalServerError)
		return
	}

	// Send presign data to the client
	resp := PresignResponse{
		FileID:         fileID,
		UploadResource: *presignedLinkResponse,
	}

	if err := RespondWithJSON(w, http.StatusOK, resp); err != nil {
		s.logger.Println("failed to send response:", err)
		return
	}
}
