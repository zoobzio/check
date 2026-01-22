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

// Validation represents the result of a single validation check.
// It tracks both the outcome (error or nil) and metadata about what was validated.
type Validation struct {
	err        error
	field      string
	validators []string
}

// Error implements the error interface.
func (v *Validation) Error() string {
	if v == nil || v.err == nil {
		return ""
	}
	return v.err.Error()
}

// Unwrap returns the underlying error for errors.Is/As compatibility.
func (v *Validation) Unwrap() error {
	if v == nil {
		return nil
	}
	return v.err
}

// Failed returns true if the validation failed.
func (v *Validation) Failed() bool {
	return v != nil && v.err != nil
}

// Result contains the aggregated outcome of multiple validations.
type Result struct {
	err     error
	applied map[string][]string
}

// Err returns the validation error (nil if validation passed).
func (r *Result) Err() error {
	if r == nil {
		return nil
	}
	return r.err
}

// Error implements the error interface for convenience.
func (r *Result) Error() string {
	if r == nil || r.err == nil {
		return ""
	}
	return r.err.Error()
}

// Unwrap implements errors.Unwrap for compatibility.
func (r *Result) Unwrap() error {
	if r == nil {
		return nil
	}
	return r.err
}

// Applied returns a map of field names to validator names that were executed.
func (r *Result) Applied() map[string][]string {
	if r == nil {
		return nil
	}
	return r.applied
}

// HasValidator checks if a specific validator was applied to a field.
func (r *Result) HasValidator(field, validator string) bool {
	if r == nil || r.applied == nil {
		return false
	}
	validators, ok := r.applied[field]
	if !ok {
		return false
	}
	for _, v := range validators {
		if v == validator {
			return true
		}
	}
	return false
}

// ValidatorsFor returns all validators applied to a specific field.
func (r *Result) ValidatorsFor(field string) []string {
	if r == nil || r.applied == nil {
		return nil
	}
	return r.applied[field]
}

// Fields returns all field names that had validators applied.
func (r *Result) Fields() []string {
	if r == nil || r.applied == nil {
		return nil
	}
	fields := make([]string, 0, len(r.applied))
	for field := range r.applied {
		fields = append(fields, field)
	}
	return fields
}

// validation creates a Validation result for a single validator.
func validation(err error, field string, validators ...string) *Validation {
	return &Validation{
		err:        err,
		field:      field,
		validators: validators,
	}
}

// All collects all validations and returns a Result.
// Tracks both successful and failed validations for metadata purposes.
func All(validations ...*Validation) *Result {
	applied := make(map[string][]string)
	var errs []error

	for _, v := range validations {
		if v == nil {
			continue
		}

		applied[v.field] = append(applied[v.field], v.validators...)

		if v.err != nil {
			errs = append(errs, v.err)
		}
	}

	var err error
	if len(errs) > 0 {
		err = Errors(errs)
	}

	return &Result{err: err, applied: applied}
}

// First returns a Result with the first failed validation, or nil error if all pass.
// Still tracks all validations that were attempted up to and including the failure.
func First(validations ...*Validation) *Result {
	applied := make(map[string][]string)

	for _, v := range validations {
		if v == nil {
			continue
		}

		applied[v.field] = append(applied[v.field], v.validators...)

		if v.err != nil {
			return &Result{err: v.err, applied: applied}
		}
	}

	return &Result{err: nil, applied: applied}
}

// Merge combines multiple Results into one.
func Merge(results ...*Result) *Result {
	applied := make(map[string][]string)
	var errs []error

	for _, r := range results {
		if r == nil {
			continue
		}
		for field, validators := range r.applied {
			applied[field] = append(applied[field], validators...)
		}
		if r.err != nil {
			var nested Errors
			if errors.As(r.err, &nested) {
				errs = append(errs, nested...)
			} else {
				errs = append(errs, r.err)
			}
		}
	}

	var err error
	if len(errs) > 0 {
		err = Errors(errs)
	}
	return &Result{err: err, applied: applied}
}

// HasErrors checks if a Result has any errors.
// Convenience function for conditional validation flows.
func HasErrors(r *Result) bool {
	return r != nil && r.err != nil
}

// GetFieldErrors extracts all FieldError values from a Result.
func GetFieldErrors(r *Result) []*FieldError {
	if r == nil || r.err == nil {
		return nil
	}
	var result []*FieldError

	var errs Errors
	if errors.As(r.err, &errs) {
		for _, e := range errs {
			var fe *FieldError
			if errors.As(e, &fe) {
				result = append(result, fe)
			}
		}
		return result
	}

	var fe *FieldError
	if errors.As(r.err, &fe) {
		result = append(result, fe)
	}
	return result
}

// FieldNames returns all field names that have errors.
func FieldNames(r *Result) []string {
	fes := GetFieldErrors(r)
	names := make([]string, 0, len(fes))
	for _, fe := range fes {
		names = append(names, fe.Field)
	}
	return names
}

// HasField checks if the Result contains a validation error for the given field.
func HasField(r *Result, field string) bool {
	for _, fe := range GetFieldErrors(r) {
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
