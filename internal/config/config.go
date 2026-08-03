package config

import (
	"github.com/spf13/viper"
)

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
	Level string `mapstructure:"level"`
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

	return cfg, nil
}
