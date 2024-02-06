package hostutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComposeAddress(t *testing.T) {
	tests := []struct {
		host        string
		port        string
		defaultPort string
		want        string
	}{
		{"", "", "8080", ":8080"},
		{"127.0.0.1", "", "8080", "127.0.0.1:8080"},
		{"127.0.0.1", "9090", "8080", "127.0.0.1:9090"},
		{"server", "9090", "8080", "server:9090"},
		{"server:1234", "9090", "8080", "server:1234"},
		{"server:1234", "", "8080", "server:1234"},
		{"server:", "9090", "8080", "server:9090"},
		{"server:", "", "8080", "server:8080"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := ComposeAddress(tt.host, tt.port, tt.defaultPort); got != tt.want {
				t.Errorf("ComposeAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComposeAddress_NoDefaultPort_Panic(t *testing.T) {
	require.Panics(t, func() {
		_ = ComposeAddress("127.0.0.1", "", "")
	})
}

func TestComposeAddresses_PortSet(t *testing.T) {
	composed := ComposeAddresses([]string{
		"",
		"127.0.0.1",
		"127.0.0.1:9090",
		"server:9090",
	}, "8080", "8080")

	require.Equal(t, []string{
		":8080",
		"127.0.0.1:8080",
		"127.0.0.1:9090",
		"server:9090",
	}, composed)
}

func TestComposeAddresses_PortNotSet(t *testing.T) {
	composed := ComposeAddresses([]string{
		"",
		"127.0.0.1",
		"127.0.0.1:9090",
		"server:9090",
	}, "", "8080")

	require.Equal(t, []string{
		":8080",
		"127.0.0.1:8080",
		"127.0.0.1:9090",
		"server:9090",
	}, composed)
}

func TestComposeAddressList_PortSet(t *testing.T) {
	composed := ComposeAddressList(", 127.0.0.1;   127.0.0.1:9090 , server:9090", "8080", "8080")

	require.Equal(t, []string{
		":8080",
		"127.0.0.1:8080",
		"127.0.0.1:9090",
		"server:9090",
	}, composed)
}

func TestComposeAddressList_PortNotSet(t *testing.T) {
	composed := ComposeAddressList(", 127.0.0.1;   127.0.0.1:9090 , server:9090", "8080", "8080")

	require.Equal(t, []string{
		":8080",
		"127.0.0.1:8080",
		"127.0.0.1:9090",
		"server:9090",
	}, composed)
}
