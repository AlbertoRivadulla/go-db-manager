package core

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

var (
	intRe = regexp.MustCompile(`^-?\d+$`)
	realRe = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	textRe = regexp.MustCompile(`^[a-zA-Z0-9 .,_-]+$`)
)

func ValidateInt(value string) error {
	if !intRe.MatchString(value) {
		return ValidationError{"value must be int"}
	}

	return nil
}

func ValidateReal(value string) error {
	if !realRe.MatchString(value) {
		return ValidationError{"value must be real"}
	}

	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return ValidationError{"could not parse into a float"}
	}

	return nil
}

func ValidateText(value string) error {
	if !textRe.MatchString(value) {
		return ValidationError{"value must be text"}
	}

	return nil
}

func ValidateTimestamp(value string) error {
	if _, err := time.Parse("2006-01-02", value); err != nil {
		return ValidationError{"value must be a date in the format YYYY-MM-DD"}
	}

	return nil
}
