// Package root provides the default configuration.
package root

import _ "embed"

// ConfigDefaults contains default configuration.
// This is the content of the .env file.
//
//go:embed .env
var ConfigDefaults string
