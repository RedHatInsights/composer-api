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

func TestLoad_InvalidLogLevel(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")

	content := []byte("log:\n  level: \"verbose\"\n")
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	chdir(t, dir)

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid log level, got nil")
	}
}

func TestLoad_EmptyPort(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")

	content := []byte("server:\n  port: \"\"\n")
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	chdir(t, dir)

	_, err := Load()
	if err == nil {
		t.Error("expected error for empty port, got nil")
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

func TestLoad_DatabaseFromFile(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")

	content := []byte(`database:
  host: "dbhost"
  port: "5432"
  name: "testdb"
  user: "admin"
  password: "secret"
  ssl_mode: "disable"
`)
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	chdir(t, dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Database.Host != "dbhost" {
		t.Errorf("expected host %q, got %q", "dbhost", cfg.Database.Host)
	}
	if cfg.Database.Port != "5432" {
		t.Errorf("expected port %q, got %q", "5432", cfg.Database.Port)
	}
	if cfg.Database.Name != "testdb" {
		t.Errorf("expected name %q, got %q", "testdb", cfg.Database.Name)
	}
	if cfg.Database.User != "admin" {
		t.Errorf("expected user %q, got %q", "admin", cfg.Database.User)
	}
	if cfg.Database.Password != "secret" {
		t.Errorf("expected password %q, got %q", "secret", cfg.Database.Password)
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("expected ssl_mode %q, got %q", "disable", cfg.Database.SSLMode)
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	tests := []struct {
		name string
		cfg  DatabaseConfig
		want string
	}{
		{
			name: "standard",
			cfg: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				Name:     "mydb",
				User:     "user",
				Password: "pass",
				SSLMode:  "disable",
			},
			want: "postgres://user:pass@localhost:5432/mydb?sslmode=disable",
		},
		{
			name: "special characters in password",
			cfg: DatabaseConfig{
				Host:     "db.host",
				Port:     "5432",
				Name:     "app",
				User:     "admin",
				Password: "p@ss:word/123",
				SSLMode:  "require",
			},
			want: "postgres://admin:p%40ss%3Aword%2F123@db.host:5432/app?sslmode=require",
		},
		{
			name: "ipv6 host",
			cfg: DatabaseConfig{
				Host:     "::1",
				Port:     "5432",
				Name:     "testdb",
				User:     "user",
				Password: "pass",
				SSLMode:  "disable",
			},
			want: "postgres://user:pass@[::1]:5432/testdb?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.DSN()
			if got != tt.want {
				t.Errorf("DSN() = %q, want %q", got, tt.want)
			}
		})
	}
}
