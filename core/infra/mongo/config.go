package mongo

// Config contains the configuration for the MongoDB connection.
type Config struct {
	Address  string
	Database string
	Username string
	Password string
}
