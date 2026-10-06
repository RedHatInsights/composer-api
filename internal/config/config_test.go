package config

import (
	"os"
	"path/filepath"
	"testing"
)

func validConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:           "8000",
			AllowedOrigins: []string{"*"},
			CORSMaxAge:     3600,
			MaxBodyBytes:   1048576,
		},
		Log: LogConfig{Level: "info"},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			Name:     "composer",
			User:     "postgres",
			Password: "postgres",
			SSLMode:  "disable",
		},
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr string
	}{
		{"valid config", nil, ""},
		{"empty server port", func(c *Config) { c.Server.Port = "" }, "server.port must not be empty"},
		{"empty allowed origins", func(c *Config) { c.Server.AllowedOrigins = nil }, "server.allowed_origins must not be empty"},
		{"zero cors max age", func(c *Config) { c.Server.CORSMaxAge = 0 }, "server.cors_max_age must be positive"},
		{"negative max body bytes", func(c *Config) { c.Server.MaxBodyBytes = -1 }, "server.max_body_bytes must be positive"},
		{"invalid log level", func(c *Config) { c.Log.Level = "verbose" }, "log.level must be one of"},
		{"empty database host", func(c *Config) { c.Database.Host = "" }, "database.host must not be empty"},
		{"empty database port", func(c *Config) { c.Database.Port = "" }, "database.port must not be empty"},
		{"empty database name", func(c *Config) { c.Database.Name = "" }, "database.name must not be empty"},
		{"empty database user", func(c *Config) { c.Database.User = "" }, "database.user must not be empty"},
		{"empty database password", func(c *Config) { c.Database.Password = "" }, "database.password must not be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			if tt.modify != nil {
				tt.modify(&cfg)
			}
			err := validate(cfg)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if got := err.Error(); len(got) < len(tt.wantErr) || got[:len(tt.wantErr)] != tt.wantErr {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, got)
			}
		})
	}
}

func writeConfig(t *testing.T, content string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
}

const fullYAML = `server:
  port: "9090"
  allowed_origins: ["*"]
  cors_max_age: 3600
  max_body_bytes: 1048576
log:
  level: "debug"
database:
  host: "dbhost"
  port: "5432"
  name: "testdb"
  user: "admin"
  password: "secret"
  ssl_mode: "require"
`

func TestLoad_FromFile(t *testing.T) {
	writeConfig(t, fullYAML)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("expected port %q, got %q", "9090", cfg.Server.Port)
	}
	if cfg.Database.Host != "dbhost" {
		t.Errorf("expected database host %q, got %q", "dbhost", cfg.Database.Host)
	}
}

func TestLoad_MissingDatabase(t *testing.T) {
	writeConfig(t, "server:\n  port: \"8000\"\nlog:\n  level: \"info\"\n")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when database config is missing")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	writeConfig(t, "invalid: [yaml: bad")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	tests := []struct {
		name string
		cfg  DatabaseConfig
		want string
	}{
		{
			"standard",
			DatabaseConfig{Host: "localhost", Port: "5432", Name: "mydb", User: "user", Password: "pass", SSLMode: "disable"},
			"postgres://user:pass@localhost:5432/mydb?sslmode=disable",
		},
		{
			"special characters in password",
			DatabaseConfig{Host: "db.host", Port: "5432", Name: "app", User: "admin", Password: "p@ss:word/123", SSLMode: "require"},
			"postgres://admin:p%40ss%3Aword%2F123@db.host:5432/app?sslmode=require",
		},
		{
			"ipv6 host",
			DatabaseConfig{Host: "::1", Port: "5432", Name: "testdb", User: "user", Password: "pass", SSLMode: "disable"},
			"postgres://user:pass@[::1]:5432/testdb?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.DSN(); got != tt.want {
				t.Errorf("DSN() = %q, want %q", got, tt.want)
			}
		})
	}
}
