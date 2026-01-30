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

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, s.Logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	file_id := q.Get("id")

	if len(file_id) == 0 {
		RespondWithError(w,
			s.Logger,
			"Missing file id in query parameter",
			errors.New("file id missing in query"),
			http.StatusBadRequest)
		return
	}

	// Get s3 key of file and delete it first from the bucket
	ObjectKey, err := s.Store.Queries.GetS3KeyFromFileID(
		r.Context(),
		database.GetS3KeyFromFileIDParams{
			UserID:  userID,
			Column2: file_id + "%",
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, s.Logger, "File not found", err, http.StatusNotFound)
			return
		}
		RespondWithError(w, s.Logger, "Error fetching s3 key", err, http.StatusInternalServerError)
		return
	}

	if len(ObjectKey) > 1 {
		RespondWithError(
			w,
			s.Logger,
			"File ID is ambiguous",
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
		RespondWithError(w, s.Logger, "Error deleting file from s3", err, http.StatusInternalServerError)
		return
	}

	err = s.Store.Queries.DeleteFileFromId(
		r.Context(),
		database.DeleteFileFromIdParams{
			UserID:  userID,
			Column2: file_id + "%",
		},
	)
	if err != nil {
		RespondWithError(w, s.Logger, "Error deleting file", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
