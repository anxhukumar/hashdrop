package handlers

import (
	"fmt"
	"net/http"
	"path"

	cloudfrontguard "github.com/anxhukumar/hashdrop/server/internal/cloudfront_guard"
	"github.com/google/uuid"
)

// Generates a signed cloudfront link and redirects to it after validating if the daily downloads limit is not hit
func (s *Server) HandlerGenerateDownloadLink(w http.ResponseWriter, r *http.Request) {

	userIDHashString := r.PathValue("userIDHash")
	fileIDString := r.PathValue("fileID")

	// Parse fileID from path
	fileID, err := uuid.Parse(fileIDString)
	if err != nil {
		RespondWithError(w, s.Logger, "Invalid fileID", err, http.StatusBadRequest)
		return
	}

	// Validate if the download attempts for this file is within the daily allowed limits
	allowed, err := cloudfrontguard.ValidateDownloadAttempts(r.Context(), s.Store.Queries, s.Cfg.DailyPerFileDownloadLimit, fileID)
	if err != nil {
		RespondWithError(w, s.Logger, "Error validating downloads attempts of the file", err, http.StatusInternalServerError)
		return
	}

	// If we are past the daily limit then give error
	if !allowed {
		RespondWithError(
			w,
			s.Logger,
			"Daily download limit of this file is exhausted",
			fmt.Errorf("Too many download requests for fileID: %s", fileIDString),
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
		RespondWithError(w, s.Logger, "Error generating signed URL", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, signedURL, http.StatusSeeOther)
}
