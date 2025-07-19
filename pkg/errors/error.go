// Package errors provides standardized error types and utilities
// for consistent error handling across the application.
//
// The package wraps and extends functionality from the standard library
// errors package, providing additional convenience functions for error wrapping.
package errors

import (
	"errors"
	"fmt"
)

// Standard error behavior interfaces from the errors package.
var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(msg string) error {
	return errors.New(msg)
}

// Wrap returns an error that wraps err with additional context.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf returns an error that wraps err with formatted additional context.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}
