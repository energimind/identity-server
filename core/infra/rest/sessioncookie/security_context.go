package sessioncookie

import (
	"fmt"
	"net"
	"net/http"
)

type securityContext struct {
	domain string
	secure bool
}

func getSecurityContext(r *http.Request) (securityContext, error) {
	domain, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		return securityContext{}, fmt.Errorf("split host port error: %w", err)
	}

	return securityContext{
		domain: domain,
		secure: r.URL.Scheme == "https",
	}, nil
}
