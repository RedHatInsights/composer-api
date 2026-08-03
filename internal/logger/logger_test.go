package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestInit_SetsDefaultLogger(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, "info", false)

	slog.Info("test message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Fatal("expected log output, got empty string")
	}
}

func TestInit_RespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, "error", false)

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

func TestInit_PrettyMode(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, "info", true)

	slog.Info("pretty test", "key", "value")

	output := buf.String()
	if output == "" {
		t.Fatal("expected log output, got empty string")
	}
	if strings.Contains(output, `"msg"`) {
		t.Error("pretty mode should not produce JSON output")
	}
}

func TestInit_JSONMode(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, "info", false)

	slog.Info("json test", "key", "value")

	output := buf.String()
	if !strings.Contains(output, `"msg"`) {
		t.Error("JSON mode should produce JSON output with \"msg\" field")
	}
}
