package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("expected default port %q, got %q", "8080", cfg.Server.Port)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("expected default log level %q, got %q", "info", cfg.Log.Level)
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")

	content := []byte("server:\n  port: \"9090\"\nlog:\n  level: \"debug\"\n")
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	chdir(t, dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("expected port %q, got %q", "9090", cfg.Server.Port)
	}

	if cfg.Log.Level != "debug" {
		t.Errorf("expected log level %q, got %q", "debug", cfg.Log.Level)
	}
}

func TestLoad_PartialConfig(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")

	content := []byte("server:\n  port: \"3000\"\n")
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	chdir(t, dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != "3000" {
		t.Errorf("expected port %q, got %q", "3000", cfg.Server.Port)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("expected default log level %q, got %q", "info", cfg.Log.Level)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")

	content := []byte("invalid: [yaml: bad")
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	chdir(t, dir)

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}
