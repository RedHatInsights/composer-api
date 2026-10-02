package config

import (
	"fmt"
	"net"
	"net/url"
	"slices"

	"github.com/spf13/viper"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

var validLogLevels = []string{"debug", "info", "warn", "error"}

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port           string   `mapstructure:"port"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	CORSMaxAge     int      `mapstructure:"cors_max_age"`
	MaxBodyBytes   int64    `mapstructure:"max_body_bytes"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func (c DatabaseConfig) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     c.Name,
		RawQuery: fmt.Sprintf("sslmode=%s", c.SSLMode),
	}
	return u.String()
}

// Load builds the application Config by merging three sources (highest precedence wins):
//  1. Clowder (cdappconfig.json) — when running on Clowder, the web provider's PublicPort and
//     database credentials are injected via v.Set() so they always override file values.
//  2. Config file (config.yaml) — searched in ".", "configs/", and "/etc/composer-api/".
//  3. Built-in defaults — fallback values used for local development when no file is present.
func Load() (Config, error) {
	v := viper.New()
	clowderCfg := clowder.LoadedConfig

	// Fallback defaults for local development (lowest precedence).
	v.SetDefault("server.port", "8000")
	v.SetDefault("server.allowed_origins", []string{"*"})
	v.SetDefault("server.cors_max_age", 3600)
	v.SetDefault("server.max_body_bytes", 1048576)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.pretty", false)

	// When deployed on Clowder, override port and database from cdappconfig.json.
	// v.Set() is used instead of v.SetDefault() so these values take precedence
	// over anything in the config file, ensuring the app binds to the port that
	// Clowder's web provider assigns and connects to the provisioned database.
	if clowder.IsClowderEnabled() {
		if clowderCfg.PublicPort != nil {
			v.Set("server.port", fmt.Sprintf("%d", *clowderCfg.PublicPort))
		}

		if db := clowderCfg.Database; db != nil {
			v.Set("database.host", db.Hostname)
			v.Set("database.port", fmt.Sprintf("%d", db.Port))
			v.Set("database.user", db.Username)
			v.Set("database.password", db.Password)
			v.Set("database.name", db.Name)
			v.Set("database.ssl_mode", db.SslMode)
		}
	}

	// Look for config.yaml in standard paths.
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("configs")
	v.AddConfigPath("/etc/composer-api")

	// Read config file if present; ignore if not found.
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	// Unmarshal the merged config into the Config struct and validate.
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("server.port must not be empty")
	}
	if len(cfg.Server.AllowedOrigins) == 0 {
		return fmt.Errorf("server.allowed_origins must not be empty")
	}
	if cfg.Server.CORSMaxAge <= 0 {
		return fmt.Errorf("server.cors_max_age must be positive")
	}
	if cfg.Server.MaxBodyBytes <= 0 {
		return fmt.Errorf("server.max_body_bytes must be positive")
	}
	if !slices.Contains(validLogLevels, cfg.Log.Level) {
		return fmt.Errorf("log.level must be one of %v, got %q", validLogLevels, cfg.Log.Level)
	}
	return nil
}
