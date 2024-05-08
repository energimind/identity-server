// Package main provides the entrypoint for the identity service.
package main

import (
	"fmt"
	"os"

	"github.com/energimind/identity-server/core/config"
	"github.com/energimind/identity-server/server"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := runCLI(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "server crashed: %v\n", err)

		os.Exit(1)
	}
}

//nolint:wrapcheck
func runCLI() error {
	cli.VersionPrinter = func(c *cli.Context) {
		_, _ = fmt.Fprintln(c.App.Writer, server.Version.Signature)
	}

	app := &cli.App{
		Name:      server.Version.Name,
		Usage:     server.Version.DisplayName,
		Version:   server.Version.Version,
		Copyright: server.Version.Copyright,
		Action: func(_ *cli.Context) error {
			return run()
		},
	}

	return app.Run(os.Args)
}

//nolint:wrapcheck
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := server.Run(cfg); err != nil {
		return err
	}

	return nil
}
