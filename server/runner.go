package server

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/energimind/go-kit/httpd"
	"github.com/energimind/go-kit/slog"
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

	slog.Info().Str("address", srv.Address()).Msg("Server listening")
	slog.Info().Msgf("%s open for business", Version.Signature)

	select {
	case err := <-errorCh:
		return err
	case sig := <-interrupted():
		slog.Info().Any("signal", sig.String()).Msg("Server interrupted")

		if err := srv.Unbind(); err != nil {
			slog.Warn().Err(err).Msg("Failed to unbind server")
		}

		if err := srv.Stop(); err != nil {
			slog.Warn().Err(err).Msg("Failed to stop server gracefully")
		}
	}

	return nil
}

func logRunTime() func() {
	startTime := time.Now()

	return func() {
		runTime := time.Since(startTime).Round(time.Second)

		slog.Info().Str("runTime", runTime.String()).Msg("Server stopped")
	}
}

func interrupted() <-chan os.Signal {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	return ch
}
