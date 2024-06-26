package server

import (
	"github.com/energimind/go-kit/slog"
	"github.com/energimind/identity-server/internal/config"
)

// Run runs the server.
func Run(cfg *config.Config) error {
	slog.Debug().Msgf("Loaded config:\n%+v", formatConfigs(config.Sections(cfg)))

	srv, clr, err := setupServer(cfg)
	if err != nil {
		return err
	}

	defer clr.closeAll()

	return run(srv)
}
