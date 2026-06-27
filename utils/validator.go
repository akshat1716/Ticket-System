package utils

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("Email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.New("Password is required")
	}
	if len(password) < 6 {
		return errors.New("Password must be at least 6 characters")
	}
	return nil
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("Name is required")
	}
	return nil
}

func ValidateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("Title is required")
	}
	return nil
}

func ValidateStatus(status string) error {
	switch status {
	case "open", "in_progress", "closed":
		return nil
	default:
		return errors.New("Invalid status. Allowed values: open, in_progress, closed")
	}
}
