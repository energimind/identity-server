package redis

// Config contains redis setup.
type Config struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Namespace  string
	Standalone bool
}
