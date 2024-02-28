package server

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/energimind/identity-server/core/infra/logger"
	"github.com/energimind/identity-server/pkg/httpd"
)

// Run runs the server. This method blocks until the server is stopped.
func run(srv *httpd.Server) error {
	defer logRunTime()()

	errorCh := make(chan error, 1)

	go func() {
		if err := srv.Run(); err != nil {
			errorCh <- err
		}
	}()

	logger.Info().Str("address", srv.Address()).Msg("Server listening")
	logger.Info().Msgf("%s open for business", Version.Signature)

	select {
	case err := <-errorCh:
		return err
	case sig := <-interrupted():
		logger.Info().Any("signal", sig.String()).Msg("Server interrupted")

		if err := srv.Unbind(); err != nil {
			logger.Warn().Err(err).Msg("Failed to unbind server")
		}

		if err := srv.Stop(); err != nil {
			logger.Warn().Err(err).Msg("Failed to stop server gracefully")
		}
	}

	return nil
}

func logRunTime() func() {
	startTime := time.Now()

	return func() {
		runTime := time.Since(startTime).Round(time.Second)

		logger.Info().Str("runTime", runTime.String()).Msg("Server stopped")
	}
}

func interrupted() <-chan os.Signal {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	return ch
}
