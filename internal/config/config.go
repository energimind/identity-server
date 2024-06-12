package config

// Config contains server setup.
type Config struct {
	HTTP   HTTPConfig
	Router RouterConfig
	Mongo  MongoConfig
	Auth   AuthenticatorConfig
	Cookie CookieConfig
	Redis  RedisConfig
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
	LocalAdminEnabled bool `env:"AUTH_LOCAL_ADMIN_ENABLED"`
}

// CookieConfig contains cookie setup.
type CookieConfig struct {
	Name   string `env:"COOKIE_NAME"`
	Secret string `env:"COOKIE_SECRET"`
}

// RedisConfig contains redis setup.
type RedisConfig struct {
	Host       string `env:"REDIS_HOST"`
	Port       string `env:"REDIS_PORT"`
	Username   string `env:"REDIS_USERNAME"`
	Password   string `env:"REDIS_PASSWORD"`
	Namespace  string `env:"REDIS_NAMESPACE"`
	Standalone bool   `env:"REDIS_STANDALONE"`
}
