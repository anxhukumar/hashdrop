package handlers

import (
	"net/http"
	"path"

	cloudfrontguard "github.com/anxhukumar/hashdrop/server/internal/cloudfront_guard"
	"github.com/google/uuid"
)

// Generates a signed cloudfront link and redirects to it after validating if the daily downloads limit is not hit
func (s *Server) HandlerGenerateDownloadLink(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_generate_download_link")

	userIDHashString := r.PathValue("userIDHash")
	fileIDString := r.PathValue("fileID")

	// Attach file_id to logger context
	logger = logger.With("file_id", fileIDString)

	// Parse fileID from path
	fileID, err := uuid.Parse(fileIDString)
	if err != nil {
		msgToDev := "invalid fileID format in path parameter"
		msgToClient := "invalid file id"
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

	// Validate if the download attempts for this file is within the daily allowed limits
	allowed, err := cloudfrontguard.ValidateDownloadAttempts(
		r.Context(),
		s.Store.Queries,
		s.Cfg.DailyPerFileDownloadLimit,
		fileID,
	)
	if err != nil {
		msgToDev := "error validating download attempts for file"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// If we are past the daily limit then give error
	if !allowed {
		msgToDev := "daily download limit exhausted for file"
		msgToClient := "daily download limit exhausted for this file"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			nil,
			http.StatusTooManyRequests,
		)
		return
	}

	// Generate signed url
	objectPath := path.Join(userIDHashString, fileIDString)
	signedURL, err := cloudfrontguard.GenerateSignedCloudfrontURL(
		s.Cfg.CloudfrontURLPrefix,
		objectPath,
		s.Cfg.CloudfrontKeyPairID,
		s.Cfg.CloudfrontPrivateKeyPath,
	)
	if err != nil {
		msgToDev := "error generating signed cloudfront url"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(w, r, signedURL, http.StatusSeeOther)
}
