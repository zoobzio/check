package check

import (
	"testing"
)

func TestNotEmptyMap(t *testing.T) {
	tests := []struct {
		value   map[string]int
		wantErr bool
	}{
		{map[string]int{"a": 1}, false},
		{map[string]int{}, true},
		{nil, true},
	}
	for _, tt := range tests {
		v := NotEmptyMap(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotEmptyMap(%v) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestEmptyMap(t *testing.T) {
	tests := []struct {
		value   map[string]int
		wantErr bool
	}{
		{map[string]int{}, false},
		{nil, false},
		{map[string]int{"a": 1}, true},
	}
	for _, tt := range tests {
		v := EmptyMap(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("EmptyMap(%v) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestMinKeys(t *testing.T) {
	tests := []struct {
		value   map[string]int
		min     int
		wantErr bool
	}{
		{map[string]int{"a": 1, "b": 2}, 2, false},
		{map[string]int{"a": 1}, 2, true},
	}
	for _, tt := range tests {
		v := MinKeys(tt.value, tt.min, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("MinKeys(%v, %d) failed = %v, wantErr %v", tt.value, tt.min, v.Failed(), tt.wantErr)
		}
	}
}

func TestMaxKeys(t *testing.T) {
	tests := []struct {
		value   map[string]int
		max     int
		wantErr bool
	}{
		{map[string]int{"a": 1}, 2, false},
		{map[string]int{"a": 1, "b": 2}, 2, false},
		{map[string]int{"a": 1, "b": 2, "c": 3}, 2, true},
	}
	for _, tt := range tests {
		v := MaxKeys(tt.value, tt.max, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("MaxKeys(%v, %d) failed = %v, wantErr %v", tt.value, tt.max, v.Failed(), tt.wantErr)
		}
	}
}

func TestHasKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	tests := []struct {
		key     string
		wantErr bool
	}{
		{"a", false},
		{"b", false},
		{"c", true},
	}
	for _, tt := range tests {
		v := HasKey(m, tt.key, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("HasKey(map, %q) failed = %v, wantErr %v", tt.key, v.Failed(), tt.wantErr)
		}
	}
}

func TestHasKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	tests := []struct {
		keys    []string
		wantErr bool
	}{
		{[]string{"a", "b"}, false},
		{[]string{"a", "b", "c"}, false},
		{[]string{"a", "d"}, true},
	}
	for _, tt := range tests {
		v := HasKeys(m, tt.keys, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("HasKeys(map, %v) failed = %v, wantErr %v", tt.keys, v.Failed(), tt.wantErr)
		}
	}
}

func TestHasAnyKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	tests := []struct {
		keys    []string
		wantErr bool
	}{
		{[]string{"a", "c"}, false},
		{[]string{"c", "d"}, true},
	}
	for _, tt := range tests {
		v := HasAnyKey(m, tt.keys, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("HasAnyKey(map, %v) failed = %v, wantErr %v", tt.keys, v.Failed(), tt.wantErr)
		}
	}
}

func TestNotHasKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	tests := []struct {
		key     string
		wantErr bool
	}{
		{"c", false},
		{"a", true},
	}
	for _, tt := range tests {
		v := NotHasKey(m, tt.key, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotHasKey(map, %q) failed = %v, wantErr %v", tt.key, v.Failed(), tt.wantErr)
		}
	}
}

func TestOnlyKeys(t *testing.T) {
	tests := []struct {
		value   map[string]int
		allowed []string
		wantErr bool
	}{
		{map[string]int{"a": 1, "b": 2}, []string{"a", "b", "c"}, false},
		{map[string]int{"a": 1, "d": 2}, []string{"a", "b", "c"}, true},
	}
	for _, tt := range tests {
		v := OnlyKeys(tt.value, tt.allowed, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("OnlyKeys(%v, %v) failed = %v, wantErr %v", tt.value, tt.allowed, v.Failed(), tt.wantErr)
		}
	}
}

func TestUniqueValues(t *testing.T) {
	tests := []struct {
		value   map[string]int
		wantErr bool
	}{
		{map[string]int{"a": 1, "b": 2}, false},
		{map[string]int{"a": 1, "b": 1}, true},
	}
	for _, tt := range tests {
		v := UniqueValues(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("UniqueValues(%v) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestEachKey(t *testing.T) {
	t.Run("all pass", func(t *testing.T) {
		m := map[string]int{"abc": 1, "def": 2}
		r := EachKey(m, func(k string) *Validation {
			return MinLen(k, 3, "key")
		})
		if r.Err() != nil {
			t.Errorf("EachKey expected nil error, got %v", r.Err())
		}
	})

	t.Run("some fail", func(t *testing.T) {
		m := map[string]int{"abc": 1, "de": 2}
		r := EachKey(m, func(k string) *Validation {
			return MinLen(k, 3, "key")
		})
		if r.Err() == nil {
			t.Error("EachKey expected error, got nil")
		}
	})

	t.Run("tracks validators", func(t *testing.T) {
		m := map[string]int{"a": 1}
		r := EachKey(m, func(k string) *Validation {
			return Required(k, "key")
		})
		if !r.HasValidator("key", "required") {
			t.Error("should track validator for key")
		}
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]int{}
		r := EachKey(m, func(k string) *Validation {
			return Required(k, "key")
		})
		if r.Err() != nil {
			t.Error("empty map should have no errors")
		}
	})
}

func TestEachMapValue(t *testing.T) {
	t.Run("all pass", func(t *testing.T) {
		m := map[string]int{"a": 10, "b": 20}
		r := EachMapValue(m, func(v int) *Validation {
			return Min(v, 5, "value")
		})
		if r.Err() != nil {
			t.Errorf("EachMapValue expected nil error, got %v", r.Err())
		}
	})

	t.Run("some fail", func(t *testing.T) {
		m := map[string]int{"a": 10, "b": 2}
		r := EachMapValue(m, func(v int) *Validation {
			return Min(v, 5, "value")
		})
		if r.Err() == nil {
			t.Error("EachMapValue expected error, got nil")
		}
	})

	t.Run("tracks validators", func(t *testing.T) {
		m := map[string]int{"a": 1}
		r := EachMapValue(m, func(v int) *Validation {
			return Positive(v, "value")
		})
		if !r.HasValidator("value", "gt") {
			t.Error("should track validator for value")
		}
	})
}

func TestEachEntry(t *testing.T) {
	t.Run("all pass", func(t *testing.T) {
		m := map[string]int{"abc": 10, "def": 20}
		r := EachEntry(m, func(k string, v int) *Validation {
			if len(k) < 3 || v < 5 {
				return validation(fieldErr("entry", "invalid"), "entry", "custom")
			}
			return validation(nil, "entry", "custom")
		})
		if r.Err() != nil {
			t.Errorf("EachEntry expected nil error, got %v", r.Err())
		}
	})

	t.Run("some fail", func(t *testing.T) {
		m := map[string]int{"ab": 10, "def": 20}
		r := EachEntry(m, func(k string, _ int) *Validation {
			if len(k) < 3 {
				return validation(fieldErr("entry", "key too short"), "entry", "custom")
			}
			return validation(nil, "entry", "custom")
		})
		if r.Err() == nil {
			t.Error("EachEntry expected error, got nil")
		}
	})

	t.Run("tracks validators", func(t *testing.T) {
		m := map[string]int{"a": 1}
		r := EachEntry(m, func(_ string, _ int) *Validation {
			return validation(nil, "entry", "custom")
		})
		if !r.HasValidator("entry", "custom") {
			t.Error("should track validator for entry")
		}
	})
}
