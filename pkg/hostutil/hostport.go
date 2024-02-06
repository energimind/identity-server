package hostutil

import (
	"net"
	"strings"
)

// HostPort encapsulates host and port.
type HostPort struct {
	Host string
	Port string
}

// NewHostPort composes HostPort from a specified address, port and a default port.
// Supported formats for address are:
//   - host
//   - host:port
//   - :port
//   - host:
//
// The method inserts specified defaultPort if port is omitted in address provided.
// The method panics if defaultPort is not specified.
// If address is empty, method will return address in form of ":defaultPort".
func NewHostPort(address, port, defaultPort string) HostPort {
	if defaultPort == "" {
		panic("missing default port")
	}

	if port == "" {
		port = defaultPort
	}

	address = strings.TrimSpace(address)

	if address == "" {
		return HostPort{Port: port}
	}

	colonAt := strings.Index(address, ":")

	if colonAt == -1 {
		return HostPort{Host: address, Port: port}
	}

	if colonAt+1 == len(address) {
		return HostPort{Host: address[0:colonAt], Port: port}
	}

	return HostPort{Host: address[0:colonAt], Port: address[colonAt+1:]}
}

// String implements Stringer interface.
func (h HostPort) String() string {
	return net.JoinHostPort(h.Host, h.Port)
}
