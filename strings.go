package check

import (
	"regexp"
	"strings"
	"unicode"
)

// Required validates that a string is not empty (after trimming whitespace).
func Required(v, field string) error {
	if strings.TrimSpace(v) == "" {
		return fieldErr(field, "is required")
	}
	return nil
}

// NotBlank validates that a string is not empty or whitespace-only.
// Unlike Required, this does not trim - it checks the raw value.
func NotBlank(v, field string) error {
	if v == "" || strings.TrimSpace(v) == "" {
		return fieldErr(field, "must not be blank")
	}
	return nil
}

// MinLen validates minimum string length (in runes, not bytes).
func MinLen(v string, minLen int, field string) error {
	if len([]rune(v)) < minLen {
		return fieldErrf(field, "must be at least %d characters", minLen)
	}
	return nil
}

// MaxLen validates maximum string length (in runes, not bytes).
func MaxLen(v string, maxLen int, field string) error {
	if len([]rune(v)) > maxLen {
		return fieldErrf(field, "must be at most %d characters", maxLen)
	}
	return nil
}

// Len validates exact string length (in runes, not bytes).
func Len(v string, exact int, field string) error {
	if len([]rune(v)) != exact {
		return fieldErrf(field, "must be exactly %d characters", exact)
	}
	return nil
}

// LenBetween validates string length is within a range (inclusive).
func LenBetween(v string, minLen, maxLen int, field string) error {
	length := len([]rune(v))
	if length < minLen || length > maxLen {
		return fieldErrf(field, "must be between %d and %d characters", minLen, maxLen)
	}
	return nil
}

// Match validates that a string matches a regular expression.
func Match(v string, pattern *regexp.Regexp, field string) error {
	if !pattern.MatchString(v) {
		return fieldErrf(field, "must match pattern %s", pattern.String())
	}
	return nil
}

// NotMatch validates that a string does not match a regular expression.
func NotMatch(v string, pattern *regexp.Regexp, field string) error {
	if pattern.MatchString(v) {
		return fieldErrf(field, "must not match pattern %s", pattern.String())
	}
	return nil
}

// Prefix validates that a string starts with the given prefix.
func Prefix(v, prefix, field string) error {
	if !strings.HasPrefix(v, prefix) {
		return fieldErrf(field, "must start with %q", prefix)
	}
	return nil
}

// Suffix validates that a string ends with the given suffix.
func Suffix(v, suffix, field string) error {
	if !strings.HasSuffix(v, suffix) {
		return fieldErrf(field, "must end with %q", suffix)
	}
	return nil
}

// Contains validates that a string contains the given substring.
func Contains(v, substr, field string) error {
	if !strings.Contains(v, substr) {
		return fieldErrf(field, "must contain %q", substr)
	}
	return nil
}

// NotContains validates that a string does not contain the given substring.
func NotContains(v, substr, field string) error {
	if strings.Contains(v, substr) {
		return fieldErrf(field, "must not contain %q", substr)
	}
	return nil
}

// OneOf validates that a string is one of the allowed values.
func OneOf(v string, allowed []string, field string) error {
	for _, a := range allowed {
		if v == a {
			return nil
		}
	}
	return fieldErrf(field, "must be one of: %s", strings.Join(allowed, ", "))
}

// NotOneOf validates that a string is not one of the disallowed values.
func NotOneOf(v string, disallowed []string, field string) error {
	for _, d := range disallowed {
		if v == d {
			return fieldErrf(field, "must not be one of: %s", strings.Join(disallowed, ", "))
		}
	}
	return nil
}

// Alpha validates that a string contains only ASCII letters.
func Alpha(v, field string) error {
	for _, r := range v {
		if !isASCIILetter(r) {
			return fieldErr(field, "must contain only letters")
		}
	}
	return nil
}

// AlphaNumeric validates that a string contains only ASCII letters and digits.
func AlphaNumeric(v, field string) error {
	for _, r := range v {
		if !isASCIILetter(r) && !isASCIIDigit(r) {
			return fieldErr(field, "must contain only letters and numbers")
		}
	}
	return nil
}

func isASCIILetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isASCIIDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Numeric validates that a string contains only ASCII digits.
func Numeric(v, field string) error {
	for _, r := range v {
		if r < '0' || r > '9' {
			return fieldErr(field, "must contain only numbers")
		}
	}
	return nil
}

// AlphaUnicode validates that a string contains only Unicode letters.
func AlphaUnicode(v, field string) error {
	for _, r := range v {
		if !unicode.IsLetter(r) {
			return fieldErr(field, "must contain only letters")
		}
	}
	return nil
}

// AlphaNumericUnicode validates that a string contains only Unicode letters and digits.
func AlphaNumericUnicode(v, field string) error {
	for _, r := range v {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fieldErr(field, "must contain only letters and numbers")
		}
	}
	return nil
}

// ASCII validates that a string contains only ASCII characters.
func ASCII(v, field string) error {
	for _, r := range v {
		if r > 127 {
			return fieldErr(field, "must contain only ASCII characters")
		}
	}
	return nil
}

// PrintableASCII validates that a string contains only printable ASCII (32-126).
func PrintableASCII(v, field string) error {
	for _, r := range v {
		if r < 32 || r > 126 {
			return fieldErr(field, "must contain only printable ASCII characters")
		}
	}
	return nil
}

// LowerCase validates that a string is entirely lowercase.
func LowerCase(v, field string) error {
	if v != strings.ToLower(v) {
		return fieldErr(field, "must be lowercase")
	}
	return nil
}

// UpperCase validates that a string is entirely uppercase.
func UpperCase(v, field string) error {
	if v != strings.ToUpper(v) {
		return fieldErr(field, "must be uppercase")
	}
	return nil
}

// NoWhitespace validates that a string contains no whitespace characters.
func NoWhitespace(v, field string) error {
	for _, r := range v {
		if unicode.IsSpace(r) {
			return fieldErr(field, "must not contain whitespace")
		}
	}
	return nil
}

// Trimmed validates that a string has no leading or trailing whitespace.
func Trimmed(v, field string) error {
	if v != strings.TrimSpace(v) {
		return fieldErr(field, "must not have leading or trailing whitespace")
	}
	return nil
}

// SingleLine validates that a string contains no newline characters.
func SingleLine(v, field string) error {
	if strings.ContainsAny(v, "\n\r") {
		return fieldErr(field, "must be a single line")
	}
	return nil
}

// Identifier validates that a string is a valid identifier (letter/underscore start, alphanumeric/underscore body).
func Identifier(v, field string) error {
	if v == "" {
		return fieldErr(field, "must not be empty")
	}
	for i, r := range v {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return fieldErr(field, "must start with a letter or underscore")
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return fieldErr(field, "must contain only letters, numbers, and underscores")
			}
		}
	}
	return nil
}

// Slug validates that a string is a valid URL slug (lowercase alphanumeric and hyphens).
func Slug(v, field string) error {
	if v == "" {
		return fieldErr(field, "must not be empty")
	}
	if v[0] == '-' || v[len(v)-1] == '-' {
		return fieldErr(field, "must not start or end with a hyphen")
	}
	for _, r := range v {
		if !isSlugChar(r) {
			return fieldErr(field, "must contain only lowercase letters, numbers, and hyphens")
		}
	}
	if strings.Contains(v, "--") {
		return fieldErr(field, "must not contain consecutive hyphens")
	}
	return nil
}

func isSlugChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-'
}
