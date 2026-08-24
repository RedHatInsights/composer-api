package config

import (
	"fmt"
	"slices"

	"github.com/spf13/viper"
)

var validLogLevels = []string{"debug", "info", "warn", "error"}

// Config holds all application configuration.
// Add new fields here as the application grows.
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
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

// Load reads configuration from a config.yaml file.
// It searches in the current directory, configs/, and /etc/composer-api/.
func Load() (Config, error) {
	v := viper.New()

	v.SetDefault("server.port", "8080")
	v.SetDefault("server.allowed_origins", []string{"*"})
	v.SetDefault("server.cors_max_age", 3600)
	v.SetDefault("server.max_body_bytes", 1048576)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.pretty", false)

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
