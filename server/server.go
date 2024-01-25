package server

import (
	"github.com/energimind/identity-service/core/config"
	"github.com/energimind/identity-service/core/infra/logger"
)

// Run runs the server.
func Run(cfg *config.Config) error {
	logger.Debug().Msgf("Loaded config:\n%+v", formatConfigs(config.Sections(cfg)))

	srv, err := buildServer(cfg)
	if err != nil {
		return err
	}

	return run(srv)
}
