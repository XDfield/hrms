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

	// APIPrefix is the prefix for all API endpoints
	APIPrefix string `envconfig:"E2E_API_PREFIX" default:""`

	// APIToken is an optional auth token
	APIToken string `envconfig:"E2E_API_TOKEN"`

	// DefaultTimeout is a general timeout for operations
	DefaultTimeout time.Duration `envconfig:"E2E_DEFAULT_TIMEOUT" default:"10s"`

	// DefaultBranchId is the default branch ID for testing
	DefaultBranchId string `envconfig:"E2E_DEFAULT_BRANCH_ID" default:"C001"`

	// AdminUser is the admin user for testing
	AdminUser string `envconfig:"E2E_ADMIN_USER" default:"admin"`

	// AdminPassword is the admin password for testing
	AdminPassword string `envconfig:"E2E_ADMIN_PASSWORD" default:"admin1"`
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