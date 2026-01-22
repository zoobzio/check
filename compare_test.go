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
		err := Equal(tt.value, tt.expected, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Equal(%d, %d) = %v, wantErr %v", tt.value, tt.expected, err, tt.wantErr)
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
		err := NotEqual(tt.value, tt.other, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotEqual(%d, %d) = %v, wantErr %v", tt.value, tt.other, err, tt.wantErr)
		}
	}
}

func TestEqualField(t *testing.T) {
	tests := []struct {
		value    string
		other    string
		wantErr  bool
	}{
		{"password123", "password123", false},
		{"password123", "password456", true},
	}
	for _, tt := range tests {
		err := EqualField(tt.value, tt.other, "password_confirm", "password")
		if (err != nil) != tt.wantErr {
			t.Errorf("EqualField(%q, %q) = %v, wantErr %v", tt.value, tt.other, err, tt.wantErr)
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
		err := NotEqualField(tt.value, tt.other, "new_password", "old_password")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotEqualField(%q, %q) = %v, wantErr %v", tt.value, tt.other, err, tt.wantErr)
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
		err := GreaterThanField(tt.value, tt.other, "end", "start")
		if (err != nil) != tt.wantErr {
			t.Errorf("GreaterThanField(%d, %d) = %v, wantErr %v", tt.value, tt.other, err, tt.wantErr)
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
		err := LessThanField(tt.value, tt.other, "start", "end")
		if (err != nil) != tt.wantErr {
			t.Errorf("LessThanField(%d, %d) = %v, wantErr %v", tt.value, tt.other, err, tt.wantErr)
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
		err := GreaterThanOrEqualField(tt.value, tt.other, "actual", "minimum")
		if (err != nil) != tt.wantErr {
			t.Errorf("GreaterThanOrEqualField(%d, %d) = %v, wantErr %v", tt.value, tt.other, err, tt.wantErr)
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
		err := LessThanOrEqualField(tt.value, tt.other, "actual", "maximum")
		if (err != nil) != tt.wantErr {
			t.Errorf("LessThanOrEqualField(%d, %d) = %v, wantErr %v", tt.value, tt.other, err, tt.wantErr)
		}
	}
}

func TestEqualStrings(t *testing.T) {
	err := Equal("hello", "hello", "field")
	if err != nil {
		t.Errorf("Equal(hello, hello) = %v, want nil", err)
	}

	err = Equal("hello", "world", "field")
	if err == nil {
		t.Error("Equal(hello, world) = nil, want error")
	}
}
