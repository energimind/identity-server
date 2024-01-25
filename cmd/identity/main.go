// Package main provides the entrypoint for the identity service.
package main

import (
	"errors"
	"slices"
	"strings"

	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/api/handler"
	"github.com/energimind/identity-service/core/config"
	"github.com/energimind/identity-service/core/infra/logger"
	"github.com/energimind/identity-service/core/infra/rest/middleware"
	"github.com/energimind/identity-service/pkg/httpd"
)

func main() {
	logger.Info().
		Str("version", "v1.0.0").
		Errs("error", []error{errors.New("some error")}).
		Msg("Starting identity service")

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Info().Msgf("Loaded config:\n%+v", formatConfigs(config.Sections(cfg)))

	router := api.NewRouter(
		middleware.LoggerInjector(),
		middleware.RequestLogger(),
		middleware.CORS(cfg.Router.AllowOrigin))

	handlers := api.Handlers{
		ApplicationHandler: handler.NewApplicationHandler(nil),
		HealthHandler:      handler.NewHealthHandler(),
	}

	routes := api.NewRoutes(handlers)

	router.RegisterRoutes(routes)

	logger.Debug().Msgf("Routes:\n%s", formatRoutes(router.GetRoutes()))

	srv, err := httpd.NewServer(httpd.Config{
		Interface: cfg.HTTP.Interface,
		Port:      cfg.HTTP.Port,
	}, router)
	if err != nil {
		panic(err)
	}

	if err := srv.Run(); err != nil {
		panic(err)
	}
}

func formatConfigs(sections []config.Section) string {
	sectionLength := 0

	for _, section := range sections {
		sectionLength = max(sectionLength, len(section.Name))
	}

	var sb strings.Builder

	for i, section := range sections {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(" -> ")
		sb.WriteString(section.Name)
		sb.WriteString(strings.Repeat(" ", sectionLength-len(section.Name)+1))
		sb.WriteRune('{')
		sb.WriteString(strings.Join(section.Values, "; "))
		sb.WriteRune('}')
	}

	return sb.String()
}

func formatRoutes(routes []api.RouteInfo) string {
	methodLength := 0

	for _, route := range routes {
		methodLength = max(methodLength, len(route.Method))
	}

	slices.SortFunc(routes, func(i1, i2 api.RouteInfo) int {
		pd := strings.Compare(i1.Path, i2.Path)

		if pd == 0 {
			// let the order be PUT, GET, DELETE
			return -strings.Compare(i1.Method, i2.Method)
		}

		return pd
	})

	var sb strings.Builder

	for i, route := range routes {
		if i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(" -> ")
		sb.WriteString(route.Method)
		sb.WriteString(strings.Repeat(" ", methodLength-len(route.Method)+1))
		sb.WriteString(route.Path)
	}

	return sb.String()
}
