package server

import "time"

// default timeouts.
const (
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 15 * time.Second
	defaultIdleTimeout       = 30 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultShutdownTimeout   = 5 * time.Second
)

// DisabledTimeout is a special value that indicates that a timeout is disabled.
// It will effectively set the timeout to a 0 value.
// It is used to override the default timeouts.
const DisabledTimeout time.Duration = -1

// timeout returns the timeout value to use. It returns 0 if the timeout is disabled.
// It returns the timeout value if the timeout is set (> 0).
// It returns the default timeout if the timeout is not set (0).
func timeout(setValue, defaultValue time.Duration) time.Duration {
	if setValue == DisabledTimeout {
		return 0
	}

	if setValue > 0 {
		return setValue
	}

	return defaultValue
}
