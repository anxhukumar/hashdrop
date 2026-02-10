package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *Server) HandlerDeleteFile(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_delete_file")

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

	q := r.URL.Query()
	fileID := q.Get("id")

	if len(fileID) == 0 {
		msgToDev := "file id missing in query parameter"
		msgToClient := "missing file id in query parameter"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			errors.New("file id missing in query"),
			http.StatusBadRequest,
		)
		return
	}

	// Attach file_id in logger context to enhance logs
	logger = logger.With("file_id", fileID)

	// Get s3 key of file and delete it first from the bucket
	ObjectKey, err := s.Store.Queries.GetS3KeyFromFileID(
		r.Context(),
		database.GetS3KeyFromFileIDParams{
			UserID:  userID,
			Column2: fileID + "%",
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msgToDev := "file not found in database for given file id"
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
		msgToDev := "error fetching s3 key of file from database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	if len(ObjectKey) > 1 {
		msgToDev := "file id prefix is ambiguous, multiple files match"
		msgToClient := "file id is ambiguous"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			errors.New("multiple files match this ID prefix"),
			http.StatusConflict,
		)
		return
	}

	// Delete object from s3
	_, err = s.S3Client.DeleteObject(r.Context(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.Cfg.S3Bucket),
		Key:    aws.String(ObjectKey[0]),
	})
	if err != nil {
		msgToDev := "error deleting file from s3"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	err = s.Store.Queries.DeleteFileFromId(
		r.Context(),
		database.DeleteFileFromIdParams{
			UserID:  userID,
			Column2: fileID + "%",
		},
	)
	if err != nil {
		msgToDev := "error deleting file record from database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
