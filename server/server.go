package server

import (
	"github.com/energimind/identity-service/core/config"
	"github.com/energimind/identity-service/core/infra/logger"
)

// Run runs the server.
func Run(cfg *config.Config) error {
	logger.Debug().Msgf("Loaded config:\n%+v", formatConfigs(config.Sections(cfg)))

	srv, clr, err := setupServer(cfg)
	if err != nil {
		clr.closeAll()

		return err
	}

	defer clr.closeAll()

	return run(srv)
}
