package config

// Config contains server setup.
type Config struct {
	HTTP   HTTPConfig
	Router RouterConfig
	Mongo  MongoConfig
	Auth   AuthenticatorConfig
	Cookie CookieConfig
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

// MongoConfig contains MongoDB setup.
type MongoConfig struct {
	Address  string `env:"MONGO_ADDRESS"`
	Database string `env:"MONGO_DATABASE"`
	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PASSWORD"`
}

// AuthenticatorConfig contains authenticator setup.
type AuthenticatorConfig struct {
	Endpoint string `env:"AUTH_ENDPOINT"`
}

// CookieConfig contains cookie setup.
type CookieConfig struct {
	Secret string
}
