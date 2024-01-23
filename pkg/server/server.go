package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Server implements an HTTP server with a complete lifecycle.
type Server struct {
	srv             *http.Server
	listener        net.Listener
	shutdownTimeout time.Duration
}

// NewServer creates a new server.
func NewServer(cfg Config, handler http.Handler) (*Server, error) {
	ln, err := net.Listen("tcp", net.JoinHostPort(cfg.Interface, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to bind HTTP server: %w", err)
	}

	return &Server{
		srv: &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: timeout(cfg.ReadHeaderTimeout, defaultReadHeaderTimeout),
			ReadTimeout:       timeout(cfg.ReadTimeout, defaultReadTimeout),
			WriteTimeout:      timeout(cfg.WriteTimeout, defaultWriteTimeout),
			IdleTimeout:       timeout(cfg.IdleTimeout, defaultIdleTimeout),
		},
		listener:        ln,
		shutdownTimeout: timeout(cfg.ShutdownTimeout, defaultShutdownTimeout),
	}, nil
}

// Run runs the server. It blocks until the server is stopped.
func (s *Server) Run() error {
	// ErrServerClosed is returned when the server is stopped.
	// ErrClosed is returned when the listener is closed.
	// We don't want to return an error in these cases.
	if err := s.srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("error serving: %w", err)
	}

	return nil
}

// Stop stops the server gracefully. It blocks until all connections are closed.
// It returns an error if the server is already stopped or if the shutdown timeout is reached.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

// Unbind unbinds the server from the listening address.
// Existing connections are not closed. New connections are rejected.
// It returns an error if the server is already unbound.
func (s *Server) Unbind() error {
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	return nil
}

// Address returns the server address.
func (s *Server) Address() string {
	return s.listener.Addr().String()
}
