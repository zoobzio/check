package check

import (
	"regexp"
	"testing"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"  hello  ", false},
		{"", true},
		{"   ", true},
		{"\t\n", true},
	}
	for _, tt := range tests {
		v := Required(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Required(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestNotBlank(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"", true},
		{"   ", true},
	}
	for _, tt := range tests {
		v := NotBlank(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotBlank(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestMinLen(t *testing.T) {
	tests := []struct {
		input   string
		min     int
		wantErr bool
	}{
		{"hello", 3, false},
		{"hi", 3, true},
		{"日本語", 3, false}, // 3 runes
		{"日本", 3, true},   // 2 runes
	}
	for _, tt := range tests {
		v := MinLen(tt.input, tt.min, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("MinLen(%q, %d) failed = %v, wantErr %v", tt.input, tt.min, v.Failed(), tt.wantErr)
		}
	}
}

func TestMaxLen(t *testing.T) {
	tests := []struct {
		input   string
		max     int
		wantErr bool
	}{
		{"hi", 3, false},
		{"hello", 3, true},
		{"日本", 3, false},  // 2 runes
		{"日本語!", 3, true}, // 4 runes
	}
	for _, tt := range tests {
		v := MaxLen(tt.input, tt.max, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("MaxLen(%q, %d) failed = %v, wantErr %v", tt.input, tt.max, v.Failed(), tt.wantErr)
		}
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		input   string
		exact   int
		wantErr bool
	}{
		{"abc", 3, false},
		{"ab", 3, true},
		{"abcd", 3, true},
	}
	for _, tt := range tests {
		v := Len(tt.input, tt.exact, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Len(%q, %d) failed = %v, wantErr %v", tt.input, tt.exact, v.Failed(), tt.wantErr)
		}
	}
}

func TestLenBetween(t *testing.T) {
	tests := []struct {
		input   string
		min     int
		max     int
		wantErr bool
	}{
		{"abc", 2, 4, false},
		{"ab", 2, 4, false},
		{"abcd", 2, 4, false},
		{"a", 2, 4, true},
		{"abcde", 2, 4, true},
	}
	for _, tt := range tests {
		v := LenBetween(tt.input, tt.min, tt.max, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("LenBetween(%q, %d, %d) failed = %v, wantErr %v", tt.input, tt.min, tt.max, v.Failed(), tt.wantErr)
		}
	}
}

func TestMatch(t *testing.T) {
	pattern := regexp.MustCompile(`^\d+$`)
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"123", false},
		{"abc", true},
		{"12a", true},
	}
	for _, tt := range tests {
		v := Match(tt.input, pattern, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Match(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestNotMatch(t *testing.T) {
	pattern := regexp.MustCompile(`secret`)
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"my secret", true},
	}
	for _, tt := range tests {
		v := NotMatch(tt.input, pattern, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotMatch(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestPrefix(t *testing.T) {
	tests := []struct {
		input   string
		prefix  string
		wantErr bool
	}{
		{"hello", "he", false},
		{"hello", "lo", true},
	}
	for _, tt := range tests {
		v := Prefix(tt.input, tt.prefix, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Prefix(%q, %q) failed = %v, wantErr %v", tt.input, tt.prefix, v.Failed(), tt.wantErr)
		}
	}
}

func TestSuffix(t *testing.T) {
	tests := []struct {
		input   string
		suffix  string
		wantErr bool
	}{
		{"hello", "lo", false},
		{"hello", "he", true},
	}
	for _, tt := range tests {
		v := Suffix(tt.input, tt.suffix, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Suffix(%q, %q) failed = %v, wantErr %v", tt.input, tt.suffix, v.Failed(), tt.wantErr)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		input   string
		substr  string
		wantErr bool
	}{
		{"hello world", "wor", false},
		{"hello world", "xyz", true},
	}
	for _, tt := range tests {
		v := Contains(tt.input, tt.substr, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Contains(%q, %q) failed = %v, wantErr %v", tt.input, tt.substr, v.Failed(), tt.wantErr)
		}
	}
}

func TestNotContains(t *testing.T) {
	tests := []struct {
		input   string
		substr  string
		wantErr bool
	}{
		{"hello world", "xyz", false},
		{"hello world", "wor", true},
	}
	for _, tt := range tests {
		v := NotContains(tt.input, tt.substr, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotContains(%q, %q) failed = %v, wantErr %v", tt.input, tt.substr, v.Failed(), tt.wantErr)
		}
	}
}

func TestOneOf(t *testing.T) {
	allowed := []string{"red", "green", "blue"}
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"red", false},
		{"green", false},
		{"yellow", true},
	}
	for _, tt := range tests {
		v := OneOf(tt.input, allowed, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("OneOf(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestNotOneOf(t *testing.T) {
	disallowed := []string{"admin", "root"}
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"user", false},
		{"admin", true},
	}
	for _, tt := range tests {
		v := NotOneOf(tt.input, disallowed, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotOneOf(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestAlpha(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"Hello", false},
		{"hello123", true},
		{"hello!", true},
	}
	for _, tt := range tests {
		v := Alpha(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Alpha(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestAlphaNumeric(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello123", false},
		{"Hello", false},
		{"hello!", true},
	}
	for _, tt := range tests {
		v := AlphaNumeric(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("AlphaNumeric(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestNumeric(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"123", false},
		{"123a", true},
		{"12.3", true},
	}
	for _, tt := range tests {
		v := Numeric(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Numeric(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestASCII(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"hello\x00", false}, // null is still ASCII
		{"héllo", true},
		{"日本語", true},
	}
	for _, tt := range tests {
		v := ASCII(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("ASCII(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestPrintableASCII(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"hello\x00", true}, // null is not printable
		{"hello\n", true},   // newline is not printable
	}
	for _, tt := range tests {
		v := PrintableASCII(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("PrintableASCII(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestLowerCase(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"Hello", true},
		{"HELLO", true},
	}
	for _, tt := range tests {
		v := LowerCase(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("LowerCase(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestUpperCase(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"HELLO", false},
		{"Hello", true},
		{"hello", true},
	}
	for _, tt := range tests {
		v := UpperCase(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("UpperCase(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestNoWhitespace(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"hello world", true},
		{"hello\t", true},
	}
	for _, tt := range tests {
		v := NoWhitespace(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NoWhitespace(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestTrimmed(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{" hello", true},
		{"hello ", true},
		{" hello ", true},
	}
	for _, tt := range tests {
		v := Trimmed(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Trimmed(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestSingleLine(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello world", false},
		{"hello\nworld", true},
		{"hello\rworld", true},
	}
	for _, tt := range tests {
		v := SingleLine(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("SingleLine(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestIdentifier(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"foo", false},
		{"_foo", false},
		{"foo_bar", false},
		{"foo123", false},
		{"123foo", true},
		{"foo-bar", true},
		{"", true},
	}
	for _, tt := range tests {
		v := Identifier(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Identifier(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestSlug(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello-world", false},
		{"hello", false},
		{"hello123", false},
		{"Hello", true},       // uppercase
		{"-hello", true},      // starts with hyphen
		{"hello-", true},      // ends with hyphen
		{"hello--world", true}, // consecutive hyphens
		{"", true},
	}
	for _, tt := range tests {
		v := Slug(tt.input, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Slug(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}
