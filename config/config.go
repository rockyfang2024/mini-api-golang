package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// UploadConfig holds file-upload-related configuration.
type UploadConfig struct {
	// Dir is the directory where uploaded files are stored (default: ./uploads).
	Dir string `mapstructure:"dir"`
	// MaxSizeMB is the maximum allowed upload size in megabytes (default: 2).
	MaxSizeMB int64 `mapstructure:"max_size_mb"`
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
	Upload   UploadConfig   `mapstructure:"upload"`
}

// LoadConfig reads and parses the configuration file using viper.
// It looks for app.yaml in the config directory or current directory.
// Environment variables with the MINI_API_ prefix override file values
// (e.g. MINI_API_SERVER_PORT, MINI_API_JWT_SECRET).
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults (mock values)
	v.SetDefault("server.port", 8080)
	v.SetDefault("database.path", "./mini-api.db")
	v.SetDefault("jwt.secret", "mock-jwt-secret-key")
	v.SetDefault("log.level", "info")
	v.SetDefault("upload.dir", "./uploads")
	v.SetDefault("upload.max_size_mb", 2)

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

	// Allow overriding any config value via MINI_API_* environment variables.
	// e.g. MINI_API_SERVER_PORT=9090 overrides server.port
	v.SetEnvPrefix("MINI_API")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}