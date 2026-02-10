package handlers

import (
	"log/slog"
	"net/http"
)

// Helper function to log warnings and also give appropirate response to the client
// To be used when the user needs an error message but our backend code is intact
func RespondWithWarn(w http.ResponseWriter, logger *slog.Logger, msgToDev, msgToClient string, err error, code int) {
	// Logs warn message for the developer
	if err != nil {
		logger.Warn(msgToDev, "msg_to_client", msgToClient, "err", err)
	} else {
		logger.Warn(msgToDev, "msg_to_client", msgToClient)
	}
	http.Error(w, msgToClient, code)
}

// Helper function to log error and also give appropriate response to the client
// To be used when the user needs an error but something in our code broke that needs to be fixed ASAP
func RespondWithError(w http.ResponseWriter, logger *slog.Logger, msgToDev string, err error, code int) {
	// Logs error message for the developer
	if err != nil {
		logger.Error(msgToDev, "msg_to_client", "internal server error", "err", err)
	} else {
		logger.Error(msgToDev, "msg_to_client", "internal server error")
	}

	// Sends error message to the client
	http.Error(w, "internal server error", code)
}
