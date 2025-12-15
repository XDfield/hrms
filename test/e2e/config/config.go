package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration for the E2E test suite
type Config struct {
	// APIEndpoint is the base URL of the service under test
	APIEndpoint string `envconfig:"E2E_API_ENDPOINT" default:"http://localhost:8888"`

	// APIToken is an optional auth token
	APIToken string `envconfig:"E2E_API_TOKEN"`

	// DefaultTimeout is a general timeout for operations
	DefaultTimeout time.Duration `envconfig:"E2E_DEFAULT_TIMEOUT" default:"30s"`

	// Database configuration for test data setup
	DbType string `envconfig:"E2E_DB_TYPE" default:"sqlite"`
	DbPath string `envconfig:"E2E_DB_PATH" default:"./data"`
	DbName string `envconfig:"E2E_DB_NAME" default:"hrms_C001"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// MustLoad loads configuration or panics on failure
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("Failed to load E2E config: %v", err)
	}
	return cfg
}
