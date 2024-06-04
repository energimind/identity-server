package sessioncookie

import (
	"net/http"
	"strings"
)

type securityContext struct {
	domain string
	secure bool
}

func getSecurityContext(r *http.Request) securityContext {
	return securityContext{
		domain: strings.Split(r.Host, ":")[0],
		secure: r.URL.Scheme == "https",
	}
}
