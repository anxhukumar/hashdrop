package logging

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // add minimum level of logs allowed
		AddSource: true,            // adds file:line automatically
	})

	return slog.New(handler)
}
