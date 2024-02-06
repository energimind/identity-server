package hostutil

import (
	"regexp"
)

// ComposeAddress composes host address from a specified address, port and a default port.
// Supported formats for address are:
//   - host
//   - host:port
//   - :port
//   - host:
//
// The method inserts specified defaultPort if port is omitted in address provided.
// The method panics if defaultPort is not specified.
// If address is empty, method will return address in form of ":defaultPort".
func ComposeAddress(address, port, defaultPort string) string {
	return NewHostPort(address, port, defaultPort).String()
}

// ComposeAddresses composes host addresses from a specified addresses, port and a default port.
func ComposeAddresses(addresses []string, port, defaultPort string) []string {
	composed := make([]string, len(addresses))

	for i := 0; i < len(addresses); i++ {
		composed[i] = ComposeAddress(addresses[i], port, defaultPort)
	}

	return composed
}

// ComposeAddressList composes host addresses from a specified addresses, port and a default port.
// Multiple addresses are delimited with comma or semi-column.
func ComposeAddressList(addresses string, port, defaultPort string) []string {
	return ComposeAddresses(regexp.MustCompile("[,;]").Split(addresses, -1), port, defaultPort)
}
