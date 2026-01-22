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
		err := NotEmptyMap(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotEmptyMap(%v) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := EmptyMap(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("EmptyMap(%v) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := MinKeys(tt.value, tt.min, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("MinKeys(%v, %d) = %v, wantErr %v", tt.value, tt.min, err, tt.wantErr)
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
		err := MaxKeys(tt.value, tt.max, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("MaxKeys(%v, %d) = %v, wantErr %v", tt.value, tt.max, err, tt.wantErr)
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
		err := HasKey(m, tt.key, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("HasKey(map, %q) = %v, wantErr %v", tt.key, err, tt.wantErr)
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
		err := HasKeys(m, tt.keys, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("HasKeys(map, %v) = %v, wantErr %v", tt.keys, err, tt.wantErr)
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
		err := HasAnyKey(m, tt.keys, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("HasAnyKey(map, %v) = %v, wantErr %v", tt.keys, err, tt.wantErr)
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
		err := NotHasKey(m, tt.key, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotHasKey(map, %q) = %v, wantErr %v", tt.key, err, tt.wantErr)
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
		err := OnlyKeys(tt.value, tt.allowed, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("OnlyKeys(%v, %v) = %v, wantErr %v", tt.value, tt.allowed, err, tt.wantErr)
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
		err := UniqueValues(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("UniqueValues(%v) = %v, wantErr %v", tt.value, err, tt.wantErr)
		}
	}
}
