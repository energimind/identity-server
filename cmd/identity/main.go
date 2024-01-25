// Package main provides the entrypoint for the identity service.
package main

import (
	"github.com/energimind/identity-service/core/config"
	"github.com/energimind/identity-service/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	if rErr := server.Run(cfg); rErr != nil {
		panic(rErr)
	}
}