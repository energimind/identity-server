package httpd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_timeout(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setValue      time.Duration
		defaultValue  time.Duration
		expectedValue time.Duration
	}{
		"disabled": {
			setValue:      DisabledTimeout,
			defaultValue:  10 * time.Second,
			expectedValue: 0,
		},
		"set": {
			setValue:      10 * time.Second,
			defaultValue:  5 * time.Second,
			expectedValue: 10 * time.Second,
		},
		"default": {
			setValue:      0,
			defaultValue:  5 * time.Second,
			expectedValue: 5 * time.Second,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, test.expectedValue, timeout(test.setValue, test.defaultValue))
		})
	}
}
