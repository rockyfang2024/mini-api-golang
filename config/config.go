package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config holds the application configuration
type Config struct {
	Port string `json:"port"`
	Env  string `json:"env"`
}

// LoadConfig reads and parses the configuration file using viper
func LoadConfig() (*Config, error) {
	var config Config

	// Set the file name of the configurations file
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("env")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")       // optionally look for config in the working directory

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s \n", err)
	}

	// Unmarshal the config into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v \n", err)
	}

	return &config, nil
}