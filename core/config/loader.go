package config

import (
	"fmt"

	root "github.com/energimind/identity-service"
	"github.com/energimind/identity-service/pkg/env"
	"github.com/energimind/identity-service/pkg/envconf"
)

// Load loads configuration from the environment.
//
// It loads the default .env file and the .env.local file if it exists.
// It also loads the .env.<ENV>.local file if the ENV environment variable is set.
func Load() (*Config, error) {
	if err := env.AutoLoad(root.ConfigDefaults); err != nil {
		return nil, fmt.Errorf("error loading embedded environment: %w", err)
	}

	return loadConfig()
}

func loadConfig() (*Config, error) {
	cfg := &Config{}

	if err := envconf.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing environment: %w", err)
	}

	return cfg, nil
}
