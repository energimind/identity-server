package hostutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHostPort(t *testing.T) {
	tests := []struct {
		address     string
		port        string
		defaultPort string
		want        HostPort
		wantString  string
	}{
		{"", "", "8080", HostPort{Port: "8080"}, ":8080"},
		{"127.0.0.1", "", "8080", HostPort{"127.0.0.1", "8080"}, "127.0.0.1:8080"},
		{"127.0.0.1", "9090", "8080", HostPort{"127.0.0.1", "9090"}, "127.0.0.1:9090"},
		{"server", "9090", "8080", HostPort{"server", "9090"}, "server:9090"},
		{"server:1234", "9090", "8080", HostPort{"server", "1234"}, "server:1234"},
		{"server:1234", "", "8080", HostPort{"server", "1234"}, "server:1234"},
		{"server:", "9090", "8080", HostPort{"server", "9090"}, "server:9090"},
		{"server:", "", "8080", HostPort{"server", "8080"}, "server:8080"},
	}

	for _, tt := range tests {
		got := NewHostPort(tt.address, tt.port, tt.defaultPort)

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("NewHostPort() = %v, want %v", got, tt.want)
		}

		if gotString := got.String(); gotString != tt.wantString {
			t.Errorf("NewHostPort() = %v, want %v", gotString, tt.wantString)
		}
	}
}

func TestNewHostPort_NoDefaultPort_Panic(t *testing.T) {
	require.Panics(t, func() {
		_ = NewHostPort("127.0.0.1", "", "")
	})
}
