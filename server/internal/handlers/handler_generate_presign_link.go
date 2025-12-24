package handlers

import (
	"net/http"
)

func (s *Server) HandlerGeneratePresignLink(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming file metadata
	var FileMetadata FileUploadRequest
	if err := DecodeJson(r, &FileMetadata); err != nil {
		RespondWithError(w, s.logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// TODO: generate presigned link with aws s3

	// TODO: upload metadata alogn with links etc to the database

	// TODO: send response data back to the client with links etc.

}
