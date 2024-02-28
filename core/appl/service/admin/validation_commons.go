package admin

import (
	"net/mail"
	"regexp"

	"github.com/energimind/identity-server/core/domain"
)

var codeRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]*$`)

func checkEmpty(name, value string) error {
	if value == "" {
		return domain.NewValidationError("%s cannot be empty", name)
	}

	return nil
}

func checkName(name string) error {
	return checkEmpty("name", name)
}

func checkCode(code string) error {
	if code == "" {
		return domain.NewValidationError("code cannot be empty")
	}

	if !codeRegex.MatchString(code) {
		return domain.NewValidationError("code contains invalid characters")
	}

	return nil
}

func checkEmail(email string) error {
	if err := checkEmpty("email", email); err != nil {
		return err
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return domain.NewValidationError("invalid email format")
	}

	return nil
}
