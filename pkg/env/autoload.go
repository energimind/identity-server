package env

import "os"

// AutoLoad loads the default .env file and the .env.local file if it exists.
//
// It also loads the .env.<ENV>.local file if the ENV environment variable is set.
func AutoLoad(defaultEnvContent string) error {
	envLoader := NewLoader()

	// embedded
	if err := envLoader.Apply(defaultEnvContent); err != nil {
		return err
	}

	// local development via selector
	selector := os.Getenv("ENV")
	if selector == "" {
		selector = "dev"
	}

	if err := envLoader.LoadOptional(".env." + selector + ".local"); err != nil {
		return err
	}

	// production
	if err := envLoader.LoadOptional(".env.local"); err != nil {
		return err
	}

	return nil
}
