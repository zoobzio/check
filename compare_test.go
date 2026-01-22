package check

import (
	"testing"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		value    int
		expected int
		wantErr  bool
	}{
		{42, 42, false},
		{42, 43, true},
	}
	for _, tt := range tests {
		v := Equal(tt.value, tt.expected, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Equal(%d, %d) failed = %v, wantErr %v", tt.value, tt.expected, v.Failed(), tt.wantErr)
		}
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		value   int
		other   int
		wantErr bool
	}{
		{42, 43, false},
		{42, 42, true},
	}
	for _, tt := range tests {
		v := NotEqual(tt.value, tt.other, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotEqual(%d, %d) failed = %v, wantErr %v", tt.value, tt.other, v.Failed(), tt.wantErr)
		}
	}
}

func TestEqualField(t *testing.T) {
	tests := []struct {
		value   string
		other   string
		wantErr bool
	}{
		{"password123", "password123", false},
		{"password123", "password456", true},
	}
	for _, tt := range tests {
		v := EqualField(tt.value, tt.other, "password_confirm", "password")
		if v.Failed() != tt.wantErr {
			t.Errorf("EqualField(%q, %q) failed = %v, wantErr %v", tt.value, tt.other, v.Failed(), tt.wantErr)
		}
	}
}

func TestNotEqualField(t *testing.T) {
	tests := []struct {
		value   string
		other   string
		wantErr bool
	}{
		{"new_pass", "old_pass", false},
		{"same_pass", "same_pass", true},
	}
	for _, tt := range tests {
		v := NotEqualField(tt.value, tt.other, "new_password", "old_password")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotEqualField(%q, %q) failed = %v, wantErr %v", tt.value, tt.other, v.Failed(), tt.wantErr)
		}
	}
}

func TestGreaterThanField(t *testing.T) {
	tests := []struct {
		value   int
		other   int
		wantErr bool
	}{
		{100, 50, false},
		{50, 50, true},
		{50, 100, true},
	}
	for _, tt := range tests {
		v := GreaterThanField(tt.value, tt.other, "end", "start")
		if v.Failed() != tt.wantErr {
			t.Errorf("GreaterThanField(%d, %d) failed = %v, wantErr %v", tt.value, tt.other, v.Failed(), tt.wantErr)
		}
	}
}

func TestLessThanField(t *testing.T) {
	tests := []struct {
		value   int
		other   int
		wantErr bool
	}{
		{50, 100, false},
		{50, 50, true},
		{100, 50, true},
	}
	for _, tt := range tests {
		v := LessThanField(tt.value, tt.other, "start", "end")
		if v.Failed() != tt.wantErr {
			t.Errorf("LessThanField(%d, %d) failed = %v, wantErr %v", tt.value, tt.other, v.Failed(), tt.wantErr)
		}
	}
}

func TestGreaterThanOrEqualField(t *testing.T) {
	tests := []struct {
		value   int
		other   int
		wantErr bool
	}{
		{100, 50, false},
		{50, 50, false},
		{49, 50, true},
	}
	for _, tt := range tests {
		v := GreaterThanOrEqualField(tt.value, tt.other, "actual", "minimum")
		if v.Failed() != tt.wantErr {
			t.Errorf("GreaterThanOrEqualField(%d, %d) failed = %v, wantErr %v", tt.value, tt.other, v.Failed(), tt.wantErr)
		}
	}
}

func TestLessThanOrEqualField(t *testing.T) {
	tests := []struct {
		value   int
		other   int
		wantErr bool
	}{
		{50, 100, false},
		{50, 50, false},
		{51, 50, true},
	}
	for _, tt := range tests {
		v := LessThanOrEqualField(tt.value, tt.other, "actual", "maximum")
		if v.Failed() != tt.wantErr {
			t.Errorf("LessThanOrEqualField(%d, %d) failed = %v, wantErr %v", tt.value, tt.other, v.Failed(), tt.wantErr)
		}
	}
}

func TestEqualStrings(t *testing.T) {
	v := Equal("hello", "hello", "field")
	if v.Failed() {
		t.Errorf("Equal(hello, hello) failed, want pass")
	}

	v = Equal("hello", "world", "field")
	if !v.Failed() {
		t.Error("Equal(hello, world) passed, want fail")
	}
}
