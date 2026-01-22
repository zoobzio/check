package check

import (
	"errors"
	"testing"
)

func TestFieldError(t *testing.T) {
	err := &FieldError{Field: "name", Message: "is required"}
	if err.Error() != "name: is required" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestErrors(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var errs Errors
		if errs.Error() != "" {
			t.Errorf("expected empty string, got: %s", errs.Error())
		}
	})

	t.Run("single", func(t *testing.T) {
		errs := Errors{&FieldError{Field: "a", Message: "bad"}}
		if errs.Error() != "a: bad" {
			t.Errorf("unexpected: %s", errs.Error())
		}
	})

	t.Run("multiple", func(t *testing.T) {
		errs := Errors{
			&FieldError{Field: "a", Message: "bad"},
			&FieldError{Field: "b", Message: "worse"},
		}
		expected := "a: bad; b: worse"
		if errs.Error() != expected {
			t.Errorf("expected %q, got %q", expected, errs.Error())
		}
	})

	t.Run("unwrap", func(t *testing.T) {
		e1 := &FieldError{Field: "a", Message: "bad"}
		errs := Errors{e1}
		unwrapped := errs.Unwrap()
		if len(unwrapped) != 1 {
			t.Error("unwrap failed: wrong length")
		}
		fe, ok := unwrapped[0].(*FieldError) //nolint:errorlint // testing identity
		if !ok || fe != e1 {
			t.Error("unwrap failed: wrong error")
		}
	})
}

func TestValidation(t *testing.T) {
	t.Run("error interface", func(t *testing.T) {
		v := &Validation{err: fieldErr("name", "is required"), field: "name", validators: []string{"required"}}
		if v.Error() != "name: is required" {
			t.Errorf("unexpected: %s", v.Error())
		}
	})

	t.Run("nil validation", func(t *testing.T) {
		var v *Validation
		if v.Error() != "" {
			t.Error("expected empty string")
		}
	})

	t.Run("passed validation", func(t *testing.T) {
		v := &Validation{err: nil, field: "name", validators: []string{"required"}}
		if v.Error() != "" {
			t.Error("expected empty string")
		}
		if v.Failed() {
			t.Error("should not be failed")
		}
	})

	t.Run("failed validation", func(t *testing.T) {
		v := &Validation{err: fieldErr("name", "is required"), field: "name", validators: []string{"required"}}
		if !v.Failed() {
			t.Error("should be failed")
		}
	})
}

func TestResult(t *testing.T) {
	t.Run("HasValidator", func(t *testing.T) {
		r := &Result{
			applied: map[string][]string{
				"email": {"required", "email"},
				"age":   {"min", "max"},
			},
		}
		if !r.HasValidator("email", "required") {
			t.Error("should have required on email")
		}
		if !r.HasValidator("email", "email") {
			t.Error("should have email on email")
		}
		if r.HasValidator("email", "min") {
			t.Error("should not have min on email")
		}
		if !r.HasValidator("age", "min") {
			t.Error("should have min on age")
		}
	})

	t.Run("Applied", func(t *testing.T) {
		r := &Result{
			applied: map[string][]string{
				"name": {"required"},
			},
		}
		applied := r.Applied()
		if len(applied["name"]) != 1 || applied["name"][0] != "required" {
			t.Errorf("unexpected applied: %v", applied)
		}
	})

	t.Run("nil Result", func(t *testing.T) {
		var r *Result
		if r.HasValidator("field", "validator") {
			t.Error("nil result should return false")
		}
		if r.Applied() != nil {
			t.Error("nil result should return nil")
		}
		if r.Err() != nil {
			t.Error("nil result should return nil error")
		}
	})

	t.Run("ValidatorsFor", func(t *testing.T) {
		r := &Result{
			applied: map[string][]string{
				"email": {"required", "email"},
			},
		}
		validators := r.ValidatorsFor("email")
		if len(validators) != 2 {
			t.Errorf("expected 2 validators, got %d", len(validators))
		}
		validators = r.ValidatorsFor("unknown")
		if validators != nil {
			t.Error("expected nil for unknown field")
		}
	})

	t.Run("Fields", func(t *testing.T) {
		r := &Result{
			applied: map[string][]string{
				"email": {"required"},
				"name":  {"required"},
			},
		}
		fields := r.Fields()
		if len(fields) != 2 {
			t.Errorf("expected 2 fields, got %d", len(fields))
		}
	})
}

func TestAll(t *testing.T) {
	t.Run("all nil", func(t *testing.T) {
		r := All(nil, nil, nil)
		if r.Err() != nil {
			t.Error("expected nil error")
		}
	})

	t.Run("collects errors", func(t *testing.T) {
		v1 := validation(fieldErr("a", "bad"), "a", "required")
		v2 := validation(fieldErr("b", "worse"), "b", "email")
		r := All(nil, v1, nil, v2, nil)
		if r.Err() == nil {
			t.Fatal("expected error")
		}
		var errs Errors
		if !errors.As(r.Err(), &errs) {
			t.Fatal("expected Errors type")
		}
		if len(errs) != 2 {
			t.Errorf("expected 2 errors, got %d", len(errs))
		}
	})

	t.Run("tracks all validators", func(t *testing.T) {
		v1 := validation(nil, "name", "required")
		v2 := validation(nil, "email", "email")
		r := All(v1, v2)
		if !r.HasValidator("name", "required") {
			t.Error("should track required on name")
		}
		if !r.HasValidator("email", "email") {
			t.Error("should track email on email")
		}
	})

	t.Run("tracks passed and failed", func(t *testing.T) {
		v1 := validation(fieldErr("name", "is required"), "name", "required")
		v2 := validation(nil, "other", "required")
		r := All(v1, v2)
		if !r.HasValidator("name", "required") {
			t.Error("should track failed validation")
		}
		if !r.HasValidator("other", "required") {
			t.Error("should track passed validation")
		}
		if r.Err() == nil {
			t.Error("should have error")
		}
	})

	t.Run("composite validators", func(t *testing.T) {
		v := validation(nil, "name", "min", "max")
		r := All(v)
		if !r.HasValidator("name", "min") {
			t.Error("should track min")
		}
		if !r.HasValidator("name", "max") {
			t.Error("should track max")
		}
	})
}

func TestFirst(t *testing.T) {
	t.Run("all nil", func(t *testing.T) {
		r := First(nil, nil)
		if r.Err() != nil {
			t.Error("expected nil error")
		}
	})

	t.Run("returns first error", func(t *testing.T) {
		v1 := validation(fieldErr("a", "first"), "a", "required")
		v2 := validation(fieldErr("b", "second"), "b", "email")
		r := First(nil, v1, v2)
		if r.Err() == nil {
			t.Fatal("expected error")
		}
		// First should stop at first error
		if !r.HasValidator("a", "required") {
			t.Error("should have tracked first validator")
		}
		// Should NOT have second validator since it stops at first error
		if r.HasValidator("b", "email") {
			t.Error("should not have tracked second validator")
		}
	})

	t.Run("tracks validators until failure", func(t *testing.T) {
		v1 := validation(nil, "name", "required") // passes
		v2 := validation(fieldErr("email", "bad"), "email", "email") // fails
		v3 := validation(nil, "age", "min") // should not be reached
		r := First(v1, v2, v3)
		if !r.HasValidator("name", "required") {
			t.Error("should have tracked first validator")
		}
		if !r.HasValidator("email", "email") {
			t.Error("should have tracked failing validator")
		}
		if r.HasValidator("age", "min") {
			t.Error("should not have tracked validator after failure")
		}
	})
}

func TestMerge(t *testing.T) {
	t.Run("merge nil", func(t *testing.T) {
		r := Merge(nil, nil)
		if r.Err() != nil {
			t.Error("expected nil error")
		}
	})

	t.Run("merge results", func(t *testing.T) {
		r1 := &Result{
			err: fieldErr("a", "bad"),
			applied: map[string][]string{
				"a": {"required"},
			},
		}
		r2 := &Result{
			err: fieldErr("b", "worse"),
			applied: map[string][]string{
				"b": {"email"},
			},
		}
		r := Merge(r1, r2)
		if r.Err() == nil {
			t.Error("expected error")
		}
		if !r.HasValidator("a", "required") {
			t.Error("should have a:required")
		}
		if !r.HasValidator("b", "email") {
			t.Error("should have b:email")
		}
	})
}

func TestHasErrors(t *testing.T) {
	if HasErrors(nil) {
		t.Error("nil should not have errors")
	}
	r := &Result{err: fieldErr("a", "b")}
	if !HasErrors(r) {
		t.Error("non-nil error should have errors")
	}
	r2 := &Result{err: nil}
	if HasErrors(r2) {
		t.Error("nil error should not have errors")
	}
}

func TestGetFieldErrors(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		fes := GetFieldErrors(nil)
		if len(fes) != 0 {
			t.Error("expected empty")
		}
	})

	t.Run("single", func(t *testing.T) {
		fe := &FieldError{Field: "a", Message: "b"}
		r := &Result{err: fe}
		fes := GetFieldErrors(r)
		if len(fes) != 1 || fes[0] != fe {
			t.Error("expected single field error")
		}
	})

	t.Run("from Errors", func(t *testing.T) {
		fe1 := &FieldError{Field: "a", Message: "1"}
		fe2 := &FieldError{Field: "b", Message: "2"}
		errs := Errors{fe1, fe2}
		r := &Result{err: errs}
		fes := GetFieldErrors(r)
		if len(fes) != 2 {
			t.Errorf("expected 2, got %d", len(fes))
		}
	})
}

func TestFieldNames(t *testing.T) {
	fe1 := &FieldError{Field: "name", Message: "1"}
	fe2 := &FieldError{Field: "email", Message: "2"}
	errs := Errors{fe1, fe2}
	r := &Result{err: errs}
	names := FieldNames(r)
	if len(names) != 2 || names[0] != "name" || names[1] != "email" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestHasField(t *testing.T) {
	fe := &FieldError{Field: "name", Message: "1"}
	errs := Errors{fe}
	r := &Result{err: errs}
	if !HasField(r, "name") {
		t.Error("should have field 'name'")
	}
	if HasField(r, "email") {
		t.Error("should not have field 'email'")
	}
}

func TestValidationTracking_Integration(t *testing.T) {
	type CreateUserInput struct {
		Email string
		Age   int
	}

	validate := func(c CreateUserInput) *Result {
		return All(
			Required(c.Email, "email"),
			Email(c.Email, "email"),
			Between(c.Age, 13, 120, "age"),
		)
	}

	t.Run("valid input", func(t *testing.T) {
		r := validate(CreateUserInput{Email: "test@example.com", Age: 25})
		if r.Err() != nil {
			t.Errorf("unexpected error: %v", r.Err())
		}
		// Verify all validators were tracked
		if !r.HasValidator("email", "required") {
			t.Error("missing required on email")
		}
		if !r.HasValidator("email", "email") {
			t.Error("missing email on email")
		}
		if !r.HasValidator("age", "min") {
			t.Error("missing min on age")
		}
		if !r.HasValidator("age", "max") {
			t.Error("missing max on age")
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		r := validate(CreateUserInput{Email: "", Age: 5})
		if r.Err() == nil {
			t.Error("expected error")
		}
		// Still tracks all validators
		if !r.HasValidator("email", "required") {
			t.Error("missing required on email")
		}
		if !r.HasValidator("email", "email") {
			t.Error("missing email on email")
		}
		if !r.HasValidator("age", "min") {
			t.Error("missing min on age")
		}
		if !r.HasValidator("age", "max") {
			t.Error("missing max on age")
		}
	})
}
