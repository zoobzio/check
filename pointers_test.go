package check

import (
	"errors"
	"testing"
)

func TestNotNil(t *testing.T) {
	val := 42
	tests := []struct {
		name    string
		value   *int
		wantErr bool
	}{
		{"not nil", &val, false},
		{"nil", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NotNil(tt.value, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("NotNil() = %v, wantErr %v", err, tt.wantErr)
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
		{"nil", nil, false},
		{"not nil", &val, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Nil(tt.value, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("Nil() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNilOr(t *testing.T) {
	t.Run("nil passes", func(t *testing.T) {
		err := NilOr[int](nil, func(_ int) error {
			return errors.New("should not be called")
		})
		if err != nil {
			t.Errorf("NilOr(nil) = %v, want nil", err)
		}
	})

	t.Run("valid value passes", func(t *testing.T) {
		val := 42
		err := NilOr(&val, func(v int) error {
			if v < 0 {
				return errors.New("must be positive")
			}
			return nil
		})
		if err != nil {
			t.Errorf("NilOr(42) = %v, want nil", err)
		}
	})

	t.Run("invalid value fails", func(t *testing.T) {
		val := -1
		err := NilOr(&val, func(v int) error {
			if v < 0 {
				return errors.New("must be positive")
			}
			return nil
		})
		if err == nil {
			t.Error("NilOr(-1) = nil, want error")
		}
	})
}

func TestNilOrField(t *testing.T) {
	t.Run("nil passes", func(t *testing.T) {
		err := NilOrField[string](nil, Required, "name")
		if err != nil {
			t.Errorf("NilOrField(nil) = %v, want nil", err)
		}
	})

	t.Run("valid value passes", func(t *testing.T) {
		val := "hello"
		err := NilOrField(&val, Required, "name")
		if err != nil {
			t.Errorf("NilOrField(hello) = %v, want nil", err)
		}
	})

	t.Run("invalid value fails", func(t *testing.T) {
		val := ""
		err := NilOrField(&val, Required, "name")
		if err == nil {
			t.Error("NilOrField('') = nil, want error")
		}
	})
}

func TestRequiredPtr(t *testing.T) {
	t.Run("nil fails", func(t *testing.T) {
		err := RequiredPtr[int](nil, func(_ int) error { return nil }, "field")
		if err == nil {
			t.Error("RequiredPtr(nil) = nil, want error")
		}
	})

	t.Run("valid value passes", func(t *testing.T) {
		val := 42
		err := RequiredPtr(&val, func(v int) error {
			if v < 0 {
				return errors.New("must be positive")
			}
			return nil
		}, "field")
		if err != nil {
			t.Errorf("RequiredPtr(42) = %v, want nil", err)
		}
	})

	t.Run("invalid value fails", func(t *testing.T) {
		val := -1
		err := RequiredPtr(&val, func(v int) error {
			if v < 0 {
				return errors.New("must be positive")
			}
			return nil
		}, "field")
		if err == nil {
			t.Error("RequiredPtr(-1) = nil, want error")
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
