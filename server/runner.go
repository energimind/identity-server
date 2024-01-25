package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/energimind/identity-service/core/infra/logger"
	"github.com/energimind/identity-service/pkg/httpd"
)

// Run runs the server. This method blocks until the server is stopped.
func run(srv *httpd.Server, releaseResources context.CancelFunc) error {
	defer logRunTime()()
	defer releaseResources()

	errorCh := make(chan error, 1)

	go func() {
		if err := srv.Run(); err != nil {
			errorCh <- err
		}
	}()

	logger.Info().Msgf("Server listening on %s", srv.Address())

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
