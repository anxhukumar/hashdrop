package handlers

import (
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
		RespondWithError(w, s.logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	file_id := q.Get("id")

	if len(file_id) == 0 {
		RespondWithError(w,
			s.logger,
			"Missing file id in query parameter",
			errors.New("file id missing in query"),
			http.StatusBadRequest)
		return
	}

	// Get s3 key of file and delete it first from the bucket
	ObjectKey, err := s.store.Queries.GetS3KeyFromFileID(
		r.Context(),
		database.GetS3KeyFromFileIDParams{
			UserID:  userID,
			Column2: file_id + "%",
		},
	)
	if err != nil {
		RespondWithError(w, s.logger, "Error fetching s3 key", err, http.StatusInternalServerError)
		return
	}

	// Delete object from s3
	_, err = s.s3Client.DeleteObject(r.Context(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.S3Bucket),
		Key:    aws.String(ObjectKey),
	})
	if err != nil {
		RespondWithError(w, s.logger, "Error deleting file from s3", err, http.StatusInternalServerError)
		return
	}

	err = s.store.Queries.DeleteFileFromId(
		r.Context(),
		database.DeleteFileFromIdParams{
			UserID:  userID,
			Column2: file_id + "%",
		},
	)
	if err != nil {
		RespondWithError(w, s.logger, "Error deleting file", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
