package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// JWTConfig holds JWT-related configuration.
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// LogConfig holds logging-related configuration.
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// Config holds the complete application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// LoadConfig reads and parses the configuration file using viper.
// It looks for app.yaml in the config directory or current directory.
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults (mock values)
	v.SetDefault("server.port", 8080)
	v.SetDefault("database.path", "./mini-api.db")
	v.SetDefault("jwt.secret", "mock-jwt-secret-key")
	v.SetDefault("log.level", "info")

	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// Read the configuration file; ignore error if file not found (defaults apply)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}