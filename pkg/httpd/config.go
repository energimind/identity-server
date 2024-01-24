package httpd

import "time"

// Config contains server setup.
type Config struct {
	Interface string
	Port      string

	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	IdleTimeout       time.Duration
	WriteTimeout      time.Duration
	ShutdownTimeout   time.Duration
}
