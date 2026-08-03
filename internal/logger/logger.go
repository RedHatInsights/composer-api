package logger

import (
	"io"
	"log/slog"
	"strings"
)

// Init creates a new slog logger and sets it as the default.
// Level should be one of "debug", "info", "warn", or "error".
func Init(destination io.Writer, level string) {
	slogLevel := parseLevel(level)
	jsonHandler := slog.NewJSONHandler(destination, &slog.HandlerOptions{AddSource: true, Level: slogLevel})
	handler := ContextHandler{Handler: jsonHandler}
	slog.SetDefault(slog.New(handler))
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
