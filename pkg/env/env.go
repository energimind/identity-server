package env

import (
	"os"
	"path"
	"strings"
)

// Loader loads environment variables from a file.
type Loader struct {
	setter func(string, string) error
}

// NewLoader creates a new Loader.
func NewLoader() *Loader {
	return &Loader{setter: os.Setenv}
}

// LoadOptional loads the environment variables from the given file.
// It ignores the file if it does not exist.
//
//nolint:wrapcheck // don't clutter the error
func (l *Loader) LoadOptional(file string) error {
	if !strings.Contains(file, "/") {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		file = path.Join(pwd, file)
	}

	content, err := os.ReadFile(path.Clean(file))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return l.apply(string(content))
}

// Apply loads the environment variables from the given content.
func (l *Loader) Apply(content string) error {
	return l.apply(content)
}

func (l *Loader) apply(content string) error {
	const expectedTokens = 2

	for _, line := range strings.Split(content, "\n") {
		if commentAt := strings.Index(line, "#"); commentAt != -1 {
			line = line[:commentAt]
		}

		if len(line) == 0 {
			continue
		}

		tokens := strings.SplitN(line, "=", expectedTokens)

		if len(tokens) != expectedTokens {
			continue
		}

		k, v := strings.TrimSpace(tokens[0]), unquote(strings.TrimSpace(tokens[1]))

		if err := l.setter(k, v); err != nil {
			return err
		}
	}

	return nil
}

func unquote(s string) string {
	const minLength = 2

	if len(s) < minLength {
		return s
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}

	if s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}

	return s
}
