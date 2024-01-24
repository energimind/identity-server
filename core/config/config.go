package config

// Config contains server setup.
type Config struct {
	HTTP   HTTPConfig
	Router RouterConfig
}

// HTTPConfig contains HTTP server setup.
type HTTPConfig struct {
	Interface string `env:"HTTPD_INTERFACE"`
	Port      string `env:"HTTPD_PORT"`
}

// RouterConfig contains router setup.
type RouterConfig struct {
	AllowOrigin string `env:"ROUTER_ALLOW_ORIGIN"`
}
