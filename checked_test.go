package check

import (
	"errors"
	"testing"
)

// Test structs with validation tags
type CheckedRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	Name     *string `json:"name" validate:"omitempty,max=100"`
	Internal string  `json:"-" validate:"-"` // Explicitly skipped
	NoTag    string  `json:"no_tag"`         // No validate tag
}

type PartiallyValidated struct {
	Required string `json:"required" validate:"required"`
	Also     string `json:"also" validate:"required"`
	Optional string `json:"optional"` // No validate tag
}

func TestChecked(t *testing.T) {
	t.Run("all tagged fields validated - passes", func(t *testing.T) {
		name := "John"
		result := Check[CheckedRequest](
			Str("test@example.com", "email").Required().Email().V(),
			Str("password123", "password").Required().MinLen(8).V(),
			OptStr(&name, "name").MaxLen(100).V(),
		)

		if result.Err() != nil {
			t.Errorf("expected pass when all tagged fields validated, got: %v", result.Err())
		}
	})

	t.Run("missing validation for tagged field - fails", func(t *testing.T) {
		// Only validate email, skip password
		result := Check[CheckedRequest](
			Str("test@example.com", "email").Required().Email().V(),
			// password not validated!
		)

		if result.Err() == nil {
			t.Error("expected error when tagged field not validated")
		}

		// Check for UncheckedFieldError
		var unchecked *UncheckedFieldError
		var errs Errors
		if !errors.As(result.Err(), &errs) {
			t.Fatalf("expected Errors type, got %T", result.Err())
		}

		found := false
		for _, err := range errs {
			if errors.As(err, &unchecked) && unchecked.Field == "password" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected UncheckedFieldError for 'password', got: %v", result.Err())
		}
	})

	t.Run("skipped field with validate:- is ignored", func(t *testing.T) {
		name := "John"
		result := Check[CheckedRequest](
			Str("test@example.com", "email").Required().Email().V(),
			Str("password123", "password").Required().MinLen(8).V(),
			OptStr(&name, "name").MaxLen(100).V(),
			// Internal field has validate:"-", should not require validation
		)

		if result.Err() != nil {
			t.Errorf("expected pass (skipped field should be ignored), got: %v", result.Err())
		}
	})

	t.Run("field without validate tag is ignored", func(t *testing.T) {
		name := "John"
		result := Check[CheckedRequest](
			Str("test@example.com", "email").Required().Email().V(),
			Str("password123", "password").Required().MinLen(8).V(),
			OptStr(&name, "name").MaxLen(100).V(),
			// no_tag field has no validate tag, should not require validation
		)

		if result.Err() != nil {
			t.Errorf("expected pass (untagged field should be ignored), got: %v", result.Err())
		}
	})

	t.Run("combines validation errors with unchecked errors", func(t *testing.T) {
		// Validate email with invalid value, skip password
		result := Check[CheckedRequest](
			Str("not-an-email", "email").Required().Email().V(),
			// password not validated!
		)

		if result.Err() == nil {
			t.Fatal("expected errors")
		}

		var errs Errors
		if !errors.As(result.Err(), &errs) {
			t.Fatalf("expected Errors type, got %T", result.Err())
		}

		// Should have both validation error and unchecked error
		hasValidationErr := false
		hasUncheckedErr := false

		for _, err := range errs {
			var fe *FieldError
			var ue *UncheckedFieldError
			if errors.As(err, &fe) {
				hasValidationErr = true
			}
			if errors.As(err, &ue) {
				hasUncheckedErr = true
			}
		}

		if !hasValidationErr {
			t.Error("expected validation error for invalid email")
		}
		if !hasUncheckedErr {
			t.Error("expected unchecked field error for password")
		}
	})

	t.Run("no validations with tagged fields - fails", func(t *testing.T) {
		result := Check[PartiallyValidated]()
		if result.Err() == nil {
			t.Error("expected error when no validations provided but fields are tagged")
		}
	})
}


func TestUncheckedFieldError(t *testing.T) {
	err := &UncheckedFieldError{
		Field:       "email",
		StructField: "Email",
		Tag:         "required,email",
	}

	expected := "email: tagged but not validated (validate: required,email)"
	if err.Error() != expected {
		t.Errorf("unexpected error message:\ngot:  %s\nwant: %s", err.Error(), expected)
	}
}

