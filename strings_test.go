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
		err := Required(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Required(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := NotBlank(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotBlank(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := MinLen(tt.input, tt.min, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("MinLen(%q, %d) = %v, wantErr %v", tt.input, tt.min, err, tt.wantErr)
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
		{"日本", 3, false}, // 2 runes
		{"日本語!", 3, true}, // 4 runes
	}
	for _, tt := range tests {
		err := MaxLen(tt.input, tt.max, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("MaxLen(%q, %d) = %v, wantErr %v", tt.input, tt.max, err, tt.wantErr)
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
		err := Len(tt.input, tt.exact, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Len(%q, %d) = %v, wantErr %v", tt.input, tt.exact, err, tt.wantErr)
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
		err := LenBetween(tt.input, tt.min, tt.max, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("LenBetween(%q, %d, %d) = %v, wantErr %v", tt.input, tt.min, tt.max, err, tt.wantErr)
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
		err := Match(tt.input, pattern, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Match(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := NotMatch(tt.input, pattern, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotMatch(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Prefix(tt.input, tt.prefix, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Prefix(%q, %q) = %v, wantErr %v", tt.input, tt.prefix, err, tt.wantErr)
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
		err := Suffix(tt.input, tt.suffix, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Suffix(%q, %q) = %v, wantErr %v", tt.input, tt.suffix, err, tt.wantErr)
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
		err := Contains(tt.input, tt.substr, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Contains(%q, %q) = %v, wantErr %v", tt.input, tt.substr, err, tt.wantErr)
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
		err := NotContains(tt.input, tt.substr, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotContains(%q, %q) = %v, wantErr %v", tt.input, tt.substr, err, tt.wantErr)
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
		err := OneOf(tt.input, allowed, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("OneOf(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := NotOneOf(tt.input, disallowed, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotOneOf(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Alpha(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Alpha(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := AlphaNumeric(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("AlphaNumeric(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Numeric(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Numeric(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := ASCII(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("ASCII(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := PrintableASCII(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("PrintableASCII(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := LowerCase(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("LowerCase(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := UpperCase(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("UpperCase(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := NoWhitespace(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NoWhitespace(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Trimmed(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Trimmed(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := SingleLine(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("SingleLine(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Identifier(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Identifier(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		{"Hello", true},      // uppercase
		{"-hello", true},     // starts with hyphen
		{"hello-", true},     // ends with hyphen
		{"hello--world", true}, // consecutive hyphens
		{"", true},
	}
	for _, tt := range tests {
		err := Slug(tt.input, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Slug(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}
