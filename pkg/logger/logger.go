package logger

import (
	"log/slog"
	"os"
)

// New creates a new structured logger configured for the Catalyst environment.
// It defaults to JSON format for production readiness.
func New(serviceName string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)

	// Add the standard service field to every log entry
	return logger.With("service", serviceName)
}
