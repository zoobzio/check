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

func TestAll(t *testing.T) {
	t.Run("all nil", func(t *testing.T) {
		err := All(nil, nil, nil)
		if err != nil {
			t.Error("expected nil")
		}
	})

	t.Run("collects errors", func(t *testing.T) {
		e1 := &FieldError{Field: "a", Message: "bad"}
		e2 := &FieldError{Field: "b", Message: "worse"}
		err := All(nil, e1, nil, e2, nil)
		var errs Errors
		if !errors.As(err, &errs) {
			t.Fatal("expected Errors type")
		}
		if len(errs) != 2 {
			t.Errorf("expected 2 errors, got %d", len(errs))
		}
	})

	t.Run("flattens nested", func(t *testing.T) {
		inner := Errors{&FieldError{Field: "a", Message: "1"}, &FieldError{Field: "b", Message: "2"}}
		outer := &FieldError{Field: "c", Message: "3"}
		err := All(inner, outer)
		var errs Errors
		if !errors.As(err, &errs) {
			t.Fatal("expected Errors type")
		}
		if len(errs) != 3 {
			t.Errorf("expected 3 errors, got %d", len(errs))
		}
	})
}

func TestFirst(t *testing.T) {
	t.Run("all nil", func(t *testing.T) {
		err := First(nil, nil)
		if err != nil {
			t.Error("expected nil")
		}
	})

	t.Run("returns first", func(t *testing.T) {
		e1 := &FieldError{Field: "a", Message: "first"}
		e2 := &FieldError{Field: "b", Message: "second"}
		err := First(nil, e1, e2)
		fe, ok := err.(*FieldError) //nolint:errorlint // testing identity
		if !ok || fe != e1 {
			t.Error("expected first error")
		}
	})
}

func TestHasErrors(t *testing.T) {
	if HasErrors(nil) {
		t.Error("nil should not have errors")
	}
	if !HasErrors(&FieldError{Field: "a", Message: "b"}) {
		t.Error("non-nil should have errors")
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
		fes := GetFieldErrors(fe)
		if len(fes) != 1 || fes[0] != fe {
			t.Error("expected single field error")
		}
	})

	t.Run("from Errors", func(t *testing.T) {
		fe1 := &FieldError{Field: "a", Message: "1"}
		fe2 := &FieldError{Field: "b", Message: "2"}
		errs := Errors{fe1, fe2}
		fes := GetFieldErrors(errs)
		if len(fes) != 2 {
			t.Errorf("expected 2, got %d", len(fes))
		}
	})
}

func TestFieldNames(t *testing.T) {
	errs := Errors{
		&FieldError{Field: "name", Message: "1"},
		&FieldError{Field: "email", Message: "2"},
	}
	names := FieldNames(errs)
	if len(names) != 2 || names[0] != "name" || names[1] != "email" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestHasField(t *testing.T) {
	errs := Errors{
		&FieldError{Field: "name", Message: "1"},
	}
	if !HasField(errs, "name") {
		t.Error("should have field 'name'")
	}
	if HasField(errs, "email") {
		t.Error("should not have field 'email'")
	}
}
