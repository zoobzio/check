package check

import (
	"testing"
)

func TestNotNil(t *testing.T) {
	val := 42
	tests := []struct {
		name    string
		value   *int
		wantErr bool
	}{
		{name: "not nil", value: &val, wantErr: false},
		{name: "nil", value: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NotNil(tt.value, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("NotNil() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
			}
		})
	}
}

func TestNil(t *testing.T) {
	val := 42
	tests := []struct {
		name    string
		value   *int
		wantErr bool
	}{
		{name: "nil", value: nil, wantErr: false},
		{name: "not nil", value: &val, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Nil(tt.value, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("Nil() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
			}
		})
	}
}

func TestNilOr(t *testing.T) {
	t.Run("nil passes", func(t *testing.T) {
		v := NilOr[int](nil, func(_ int) *Validation {
			return validation(fieldErr("field", "should not be called"), "field", "test")
		})
		if v != nil {
			t.Errorf("NilOr(nil) = %v, want nil", v)
		}
	})

	t.Run("valid value passes", func(t *testing.T) {
		val := 42
		v := NilOr(&val, func(v int) *Validation {
			if v < 0 {
				return validation(fieldErr("field", "must be positive"), "field", "positive")
			}
			return validation(nil, "field", "positive")
		})
		if v.Failed() {
			t.Errorf("NilOr(42) failed, want pass")
		}
	})

	t.Run("invalid value fails", func(t *testing.T) {
		val := -1
		v := NilOr(&val, func(v int) *Validation {
			if v < 0 {
				return validation(fieldErr("field", "must be positive"), "field", "positive")
			}
			return validation(nil, "field", "positive")
		})
		if !v.Failed() {
			t.Error("NilOr(-1) = pass, want fail")
		}
	})
}

func TestNilOrField(t *testing.T) {
	t.Run("nil passes", func(t *testing.T) {
		v := NilOrField[string](nil, Required, "name")
		if v != nil {
			t.Errorf("NilOrField(nil) = %v, want nil", v)
		}
	})

	t.Run("valid value passes", func(t *testing.T) {
		val := "hello"
		v := NilOrField(&val, Required, "name")
		if v.Failed() {
			t.Errorf("NilOrField(hello) failed, want pass")
		}
	})

	t.Run("invalid value fails", func(t *testing.T) {
		val := ""
		v := NilOrField(&val, Required, "name")
		if !v.Failed() {
			t.Error("NilOrField('') = pass, want fail")
		}
	})
}

func TestRequiredPtr(t *testing.T) {
	t.Run("nil fails", func(t *testing.T) {
		v := RequiredPtr[int](nil, func(_ int) *Validation { return nil }, "field")
		if !v.Failed() {
			t.Error("RequiredPtr(nil) = pass, want fail")
		}
		if v.validators[0] != "required" {
			t.Error("should have required validator")
		}
	})

	t.Run("valid value passes", func(t *testing.T) {
		val := 42
		v := RequiredPtr(&val, func(v int) *Validation {
			if v < 0 {
				return validation(fieldErr("field", "must be positive"), "field", "positive")
			}
			return validation(nil, "field", "positive")
		}, "field")
		if v.Failed() {
			t.Errorf("RequiredPtr(42) failed, want pass")
		}
	})

	t.Run("invalid value fails", func(t *testing.T) {
		val := -1
		v := RequiredPtr(&val, func(v int) *Validation {
			if v < 0 {
				return validation(fieldErr("field", "must be positive"), "field", "positive")
			}
			return validation(nil, "field", "positive")
		}, "field")
		if !v.Failed() {
			t.Error("RequiredPtr(-1) = pass, want fail")
		}
	})

	t.Run("combines validators", func(t *testing.T) {
		val := 42
		v := RequiredPtr(&val, func(_ int) *Validation {
			return validation(nil, "field", "positive")
		}, "field")
		// Should have both "required" and "positive"
		hasRequired := false
		hasPositive := false
		for _, name := range v.validators {
			if name == "required" {
				hasRequired = true
			}
			if name == "positive" {
				hasPositive = true
			}
		}
		if !hasRequired || !hasPositive {
			t.Errorf("expected validators [required, positive], got %v", v.validators)
		}
	})
}

func TestDeref(t *testing.T) {
	t.Run("not nil", func(t *testing.T) {
		val := 42
		result := Deref(&val)
		if result != 42 {
			t.Errorf("Deref(&42) = %d, want 42", result)
		}
	})

	t.Run("nil", func(t *testing.T) {
		result := Deref[int](nil)
		if result != 0 {
			t.Errorf("Deref(nil) = %d, want 0", result)
		}
	})
}

func TestDerefOr(t *testing.T) {
	t.Run("not nil", func(t *testing.T) {
		val := 42
		result := DerefOr(&val, 99)
		if result != 42 {
			t.Errorf("DerefOr(&42, 99) = %d, want 42", result)
		}
	})

	t.Run("nil", func(t *testing.T) {
		result := DerefOr[int](nil, 99)
		if result != 99 {
			t.Errorf("DerefOr(nil, 99) = %d, want 99", result)
		}
	})
}

func TestPtr(t *testing.T) {
	result := Ptr(42)
	if result == nil || *result != 42 {
		t.Error("Ptr(42) failed")
	}
}
