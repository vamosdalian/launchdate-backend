package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

// Config holds all configuration for the application
type Config struct {
	Server             ServerConfig
	MongodbURL         string `env:"MONGODB_URL"`
	MongodbDatabase    string `env:"MONGODB_DATABASE"`
	LL2URLPrefix       string `env:"LL2_URL_PREFIX"`
	LL2RequestInterval int    `env:"LL2_REQUEST_INTERVAL, default=5"` // in seconds
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `env:"SERVER_PORT"`
	Host string `env:"SERVER_HOST"`
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
