// Package check provides zero-reflection validation primitives for Go.
// Explicit validation functions, no struct tags, fully composable.
package check

import (
	"errors"
	"fmt"
	"strings"
)

// FieldError represents a validation error for a specific field.
type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Errors is a collection of validation errors.
type Errors []error

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	var b strings.Builder
	for i, err := range e {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

// Unwrap returns the underlying errors for use with errors.Is/As.
func (e Errors) Unwrap() []error {
	return e
}

// All collects all non-nil errors into a single Errors value.
// Returns nil if all errors are nil.
func All(errs ...error) error {
	var collected Errors
	for _, err := range errs {
		if err != nil {
			// Flatten nested Errors
			var nested Errors
			if errors.As(err, &nested) {
				collected = append(collected, nested...)
			} else {
				collected = append(collected, err)
			}
		}
	}
	if len(collected) == 0 {
		return nil
	}
	return collected
}

// First returns the first non-nil error, or nil if all are nil.
// Use this when you want fail-fast behavior instead of collecting all errors.
func First(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// HasErrors checks if an error is non-nil.
// Convenience function for conditional validation flows.
func HasErrors(err error) bool {
	return err != nil
}

// GetFieldErrors extracts all FieldError values from an error.
// Works with both single FieldError and Errors collections.
func GetFieldErrors(err error) []*FieldError {
	if err == nil {
		return nil
	}
	var result []*FieldError

	// Check if it's an Errors collection first
	var errs Errors
	if errors.As(err, &errs) {
		for _, e := range errs {
			var fe *FieldError
			if errors.As(e, &fe) {
				result = append(result, fe)
			}
		}
		return result
	}

	// Otherwise check if it's a single FieldError
	var fe *FieldError
	if errors.As(err, &fe) {
		result = append(result, fe)
	}
	return result
}

// FieldNames returns all field names that have errors.
func FieldNames(err error) []string {
	fes := GetFieldErrors(err)
	names := make([]string, 0, len(fes))
	for _, fe := range fes {
		names = append(names, fe.Field)
	}
	return names
}

// HasField checks if the error contains a validation error for the given field.
func HasField(err error, field string) bool {
	for _, fe := range GetFieldErrors(err) {
		if fe.Field == field {
			return true
		}
	}
	return false
}

// fieldErr creates a FieldError with the given field and message.
func fieldErr(field, message string) error {
	return &FieldError{Field: field, Message: message}
}

// fieldErrf creates a FieldError with a formatted message.
func fieldErrf(field, format string, args ...any) error {
	return &FieldError{Field: field, Message: fmt.Sprintf(format, args...)}
}
