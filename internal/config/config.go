package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

// Config holds all configuration for the application
type Config struct {
	Server          ServerConfig
	MongodbURL      string `env:"MONGODB_URL"`
	MongodbDatabase string `env:"MONGODB_DATABASE"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `env:"PORT"`
	Host string `env:"HOST"`
	Env  string `env:"ENVIRONMENT"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{}
	if err := envconfig.Process(context.Background(), config); err != nil {
		return nil, err
	}

	return config, nil
}
