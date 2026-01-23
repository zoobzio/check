package check

import (
	"errors"
	"strings"

	"github.com/zoobzio/sentinel"
)

// Check validates the given validations and verifies that all fields with validate tags were checked.
// This is the primary API for struct validation - it combines validation execution with tag verification.
//
// Usage:
//
//	func (r *Request) Validate() error {
//	    return check.Check[Request](
//	        check.Str(r.Email, "email").Required().Email().V(),
//	        check.Str(r.Name, "name").Required().V(),
//	    ).Err()
//	}
func Check[T any](validations ...*Validation) *Result {
	result := All(validations...)

	// Inspect the type to get field metadata
	metadata := sentinel.Inspect[T]()

	// Get the fields that were actually validated
	applied := result.Applied()
	if applied == nil {
		applied = make(map[string][]string)
	}

	// Check each field with a validate tag
	var missingErrs []error
	for _, field := range metadata.Fields {
		validateTag, hasTag := field.Tags["validate"]
		if !hasTag || validateTag == "" || validateTag == "-" {
			continue
		}

		// Determine the field name to check (prefer json tag, fall back to field name)
		fieldName := getFieldName(field)

		// Check if this field was validated using either the json tag name or struct field name
		_, validatedByFieldName := applied[fieldName]
		_, validatedByStructName := applied[field.Name]
		if !validatedByFieldName && !validatedByStructName {
			missingErrs = append(missingErrs, &UncheckedFieldError{
				Field:       fieldName,
				StructField: field.Name,
				Tag:         validateTag,
			})
		}
	}

	// If no missing validations, return original result
	if len(missingErrs) == 0 {
		return result
	}

	// Combine original errors with missing validation errors
	var allErrs Errors
	if result.err != nil {
		var existing Errors
		if errors.As(result.err, &existing) {
			allErrs = append(allErrs, existing...)
		} else {
			allErrs = append(allErrs, result.err)
		}
	}
	allErrs = append(allErrs, missingErrs...)

	return &Result{
		err:     allErrs,
		applied: applied,
	}
}

// getFieldName determines the field name used in validation.
// Prefers json tag name, falls back to lowercase struct field name.
func getFieldName(field sentinel.FieldMetadata) string {
	if jsonTag, ok := field.Tags["json"]; ok && jsonTag != "" && jsonTag != "-" {
		// Handle json tag options like `json:"name,omitempty"`
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonTag = jsonTag[:idx]
		}
		if jsonTag != "" {
			return jsonTag
		}
	}
	return strings.ToLower(field.Name)
}

// UncheckedFieldError indicates a field has validation requirements but was not validated.
type UncheckedFieldError struct {
	Field       string // The field name used in validation (json tag or lowercase)
	StructField string // The actual struct field name
	Tag         string // The validate tag value
}

func (e *UncheckedFieldError) Error() string {
	return e.Field + ": tagged but not validated (validate: " + e.Tag + ")"
}

