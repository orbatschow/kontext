package logger

import (
	"log/slog"
	"os"
)

const DefaultVerbosity = 2

var Verbosity = DefaultVerbosity

// New creates a new slog.Logger with the configured verbosity level.
// Verbosity maps directly to slog levels:
// 1=DEBUG, 2=INFO, 3=WARN, 4=ERROR
func New() *slog.Logger {
	var level slog.Level
	switch Verbosity {
	case 1:
		level = slog.LevelDebug
	case 2:
		level = slog.LevelInfo
	case 3:
		level = slog.LevelWarn
	case 4:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stderr, opts)
	return slog.New(handler)
}
