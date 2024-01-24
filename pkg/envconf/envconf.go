package envconf

import (
	"fmt"
	"reflect"
	"strconv"

	envs "github.com/caarlos0/env/v7"
)

// Parse parses the environment variables into the given struct.
//
// It uses the default parsers from caarlos0/env, and adds some custom parsers
// for bool values.
func Parse(v any) error {
	if err := envs.ParseWithFuncs(v, customParsers()); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	return nil
}

func customParsers() map[reflect.Type]envs.ParserFunc {
	return map[reflect.Type]envs.ParserFunc{
		reflect.TypeOf(true): func(v string) (any, error) {
			switch v {
			case "on", "yes":
				return true, nil
			case "off", "no":
				return false, nil
			default:
				//nolint:wrapcheck // don't clutter the error
				return strconv.ParseBool(v)
			}
		},
	}
}
