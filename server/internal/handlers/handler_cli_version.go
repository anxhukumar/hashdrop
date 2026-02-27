package handlers

import (
	"net/http"
)

// Sends the current cli version compatible with the server
func (s *Server) HandlerCliVersion(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_cli_version")

	// Cli version response
	cliVersionResponse := CliVersion{
		CompatibleVersion: s.Cfg.CliVersion,
	}

	if err := RespondWithJSON(w, http.StatusOK, cliVersionResponse); err != nil {
		logger.Error("failed to send cli version response", "err", err)
	}
}
