package server

import (
	"fmt"
	"slices"
	"strings"

	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/api/handler"
	"github.com/energimind/identity-service/core/config"
	"github.com/energimind/identity-service/core/infra/logger"
	"github.com/energimind/identity-service/core/infra/rest/middleware"
	"github.com/energimind/identity-service/core/infra/rest/router"
	"github.com/energimind/identity-service/pkg/httpd"
)

// buildServer creates and configures a new server.
func buildServer(cfg *config.Config) (*httpd.Server, error) {
	handlers := api.Handlers{
		ApplicationHandler: handler.NewApplicationHandler(nil),
		HealthHandler:      handler.NewHealthHandler(),
	}

	routes := api.NewRoutes(handlers)

	restRouter := router.New(
		middleware.LoggerInjector(),
		middleware.RequestLogger(),
		middleware.CORS(cfg.Router.AllowOrigin))

	restRouter.RegisterRoutes(routes)

	srv, err := httpd.NewServer(httpd.Config{
		Interface: cfg.HTTP.Interface,
		Port:      cfg.HTTP.Port,
	}, restRouter)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	logger.Debug().Msgf("Routes:\n%s", formatRoutes(restRouter.GetRoutes()))

	return srv, nil
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

func formatRoutes(routes []router.RouteInfo) string {
	methodLength := 0

	for _, route := range routes {
		methodLength = max(methodLength, len(route.Method))
	}

	slices.SortFunc(routes, func(i1, i2 router.RouteInfo) int {
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
