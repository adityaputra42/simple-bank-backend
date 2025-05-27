package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidateUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidateFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxlength int) error {

	n := len(value)

	if n < minLength || n > maxlength {
		return fmt.Errorf("must content from %d-%d characters", minLength, maxlength)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateUsername(value) {
		return fmt.Errorf("must contains only lowercase letters, digit or underscore")
	}
	return nil
}
func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidateFullName(value) {
		return fmt.Errorf("must contains only letters or spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 8, 100); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not valid email address")
	}
	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}

	return nil
}

func ValidateScretCode(value string) error {

	return ValidateString(value, 32, 128)
}
