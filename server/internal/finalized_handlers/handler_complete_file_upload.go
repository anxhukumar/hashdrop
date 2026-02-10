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
	logger := s.Logger.With("handler", "handler_complete_file_upload")

	// Get decoded incoming file json data
	var FileUploadMetadata FileUploadMetadata
	if err := DecodeJson(r, &FileUploadMetadata); err != nil {
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

	// Attach fileID in logger context to enhance logs
	logger = logger.With("file_id", FileUploadMetadata.FileID)

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

	// Attach UserId in logger context to enhance logs
	logger = logger.With("user_id", userID.String())

	// Get metadata of the uploaded file

	// Fetch s3ObjectKey from db
	ObjectKey, err := s.Store.Queries.GetS3KeyForUploadVerification(r.Context(), database.GetS3KeyForUploadVerificationParams{
		ID:     FileUploadMetadata.FileID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msgToDev := "the file that the user is trying to verify does not exist in database"
			msgToClient := "file not found"
			RespondWithWarn(
				w,
				logger,
				msgToDev,
				msgToClient,
				err,
				http.StatusNotFound,
			)
			return
		}
		msgToDev := "error fetching S3ObjectKey from database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	head, err := s.S3Client.HeadObject(r.Context(), &s3.HeadObjectInput{
		Bucket: aws.String(s.Cfg.S3Bucket),
		Key:    aws.String(ObjectKey),
	})
	if err != nil {
		msgToDev := "error fetching object metadata from s3"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	verifiedFileSize := aws.ToInt64(head.ContentLength) // int64
	if verifiedFileSize > s.Cfg.S3PerFileMaxDataSize {
		// Delete object
		_, _ = s.S3Client.DeleteObject(r.Context(), &s3.DeleteObjectInput{
			Bucket: aws.String(s.Cfg.S3Bucket),
			Key:    aws.String(ObjectKey),
		})

		// Update db status to failed
		if err := s.Store.Queries.UpdateFailedFile(r.Context(), database.UpdateFailedFileParams{
			Status: "failed",
			ID:     FileUploadMetadata.FileID,
		}); err != nil {
			msgToDev := "error marking file upload status as failed"
			RespondWithError(
				w,
				logger,
				msgToDev,
				err,
				http.StatusInternalServerError,
			)
			return
		}

		msgToDev := "the file that the user tried verifying is above the allowed size limit"
		msgToClient := "file size exceeds the allowed limit"
		logger = logger.With("file_size", verifiedFileSize)
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			nil,
			http.StatusRequestEntityTooLarge,
		)
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
	if err := s.Store.Queries.UpdateUploadedFile(r.Context(), fileData); err != nil {
		msgToDev := "error updating file meta data"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	RespondWithJSON(
		w,
		http.StatusOK,
		FileUploadSuccessResponse{
			S3ObjectKey:      ObjectKey,
			UploadedFileSize: verifiedFileSize,
		})

	logger.Info("uploaded file verified successfully")

}
