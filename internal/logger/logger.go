package logger

import (
	"io"
	"log/slog"
	"strings"
)

// Init creates a new slog logger and sets it as the default.
// Level should be one of "debug", "info", "warn", or "error".
// When pretty is true, logs are formatted as human-readable text for local development.
func Init(destination io.Writer, level string, pretty bool) {
	slogLevel := parseLevel(level)
	opts := &slog.HandlerOptions{AddSource: true, Level: slogLevel}

	var base slog.Handler
	if pretty {
		base = slog.NewTextHandler(destination, opts)
	} else {
		base = slog.NewJSONHandler(destination, opts)
	}

	slog.SetDefault(slog.New(ContextHandler{Handler: base}))
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
