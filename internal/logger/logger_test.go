package logger

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestInit_SetsDefaultLogger(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, "info")

	slog.Info("test message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Fatal("expected log output, got empty string")
	}
}

func TestInit_RespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, "error")

	slog.Info("should not appear")

	if buf.Len() != 0 {
		t.Errorf("expected no output for info at error level, got %q", buf.String())
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"unknown", slog.LevelInfo},
		{"", slog.LevelInfo},
	}

	for _, tt := range tests {
		if got := parseLevel(tt.input); got != tt.expected {
			t.Errorf("parseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}
