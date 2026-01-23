package check

import (
	"fmt"
	"regexp"

	"golang.org/x/exp/constraints"
)

// -----------------------------------------------------------------------------
// Internal helpers
// -----------------------------------------------------------------------------

// combine merges multiple validations for the same field into one.
func combine(field string, validations []*Validation) *Validation {
	if len(validations) == 0 {
		return nil
	}

	var errs Errors
	var validators []string

	for _, v := range validations {
		if v == nil {
			continue
		}
		validators = append(validators, v.validators...)
		if v.err != nil {
			errs = append(errs, v.err)
		}
	}

	if len(validators) == 0 {
		return nil
	}

	var err error
	if len(errs) == 1 {
		err = errs[0]
	} else if len(errs) > 1 {
		err = errs
	}

	return &Validation{err: err, field: field, validators: validators}
}

// -----------------------------------------------------------------------------
// String Builder
// -----------------------------------------------------------------------------

// StrBuilder provides fluent validation for string values.
type StrBuilder struct {
	value       string
	field       string
	validations []*Validation
}

// Str creates a new string validation builder.
func Str(v string, field string) *StrBuilder {
	return &StrBuilder{value: v, field: field}
}

// V returns the combined validation result.
func (b *StrBuilder) V() *Validation {
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *StrBuilder) When(cond bool, fn func(*StrBuilder)) *StrBuilder {
	if cond {
		fn(b)
	}
	return b
}

// Required validates that the string is not empty.
func (b *StrBuilder) Required() *StrBuilder {
	b.validations = append(b.validations, Required(b.value, b.field))
	return b
}

// NotBlank validates that the string is not empty or whitespace-only.
func (b *StrBuilder) NotBlank() *StrBuilder {
	b.validations = append(b.validations, NotBlank(b.value, b.field))
	return b
}

// MinLen validates minimum string length.
func (b *StrBuilder) MinLen(n int) *StrBuilder {
	b.validations = append(b.validations, MinLen(b.value, n, b.field))
	return b
}

// MaxLen validates maximum string length.
func (b *StrBuilder) MaxLen(n int) *StrBuilder {
	b.validations = append(b.validations, MaxLen(b.value, n, b.field))
	return b
}

// Len validates exact string length.
func (b *StrBuilder) Len(n int) *StrBuilder {
	b.validations = append(b.validations, Len(b.value, n, b.field))
	return b
}

// LenBetween validates string length is within a range.
func (b *StrBuilder) LenBetween(minLen, maxLen int) *StrBuilder {
	b.validations = append(b.validations, LenBetween(b.value, minLen, maxLen, b.field))
	return b
}

// Match validates that the string matches a pattern.
func (b *StrBuilder) Match(pattern *regexp.Regexp) *StrBuilder {
	b.validations = append(b.validations, Match(b.value, pattern, b.field))
	return b
}

// NotMatch validates that the string does not match a pattern.
func (b *StrBuilder) NotMatch(pattern *regexp.Regexp) *StrBuilder {
	b.validations = append(b.validations, NotMatch(b.value, pattern, b.field))
	return b
}

// Prefix validates that the string starts with the given prefix.
func (b *StrBuilder) Prefix(prefix string) *StrBuilder {
	b.validations = append(b.validations, Prefix(b.value, prefix, b.field))
	return b
}

// Suffix validates that the string ends with the given suffix.
func (b *StrBuilder) Suffix(suffix string) *StrBuilder {
	b.validations = append(b.validations, Suffix(b.value, suffix, b.field))
	return b
}

// Contains validates that the string contains the substring.
func (b *StrBuilder) Contains(substr string) *StrBuilder {
	b.validations = append(b.validations, Contains(b.value, substr, b.field))
	return b
}

// NotContains validates that the string does not contain the substring.
func (b *StrBuilder) NotContains(substr string) *StrBuilder {
	b.validations = append(b.validations, NotContains(b.value, substr, b.field))
	return b
}

// OneOf validates that the string is one of the allowed values.
func (b *StrBuilder) OneOf(allowed []string) *StrBuilder {
	b.validations = append(b.validations, OneOf(b.value, allowed, b.field))
	return b
}

// NotOneOf validates that the string is not one of the disallowed values.
func (b *StrBuilder) NotOneOf(disallowed []string) *StrBuilder {
	b.validations = append(b.validations, NotOneOf(b.value, disallowed, b.field))
	return b
}

// Alpha validates that the string contains only ASCII letters.
func (b *StrBuilder) Alpha() *StrBuilder {
	b.validations = append(b.validations, Alpha(b.value, b.field))
	return b
}

// AlphaNumeric validates that the string contains only ASCII letters and digits.
func (b *StrBuilder) AlphaNumeric() *StrBuilder {
	b.validations = append(b.validations, AlphaNumeric(b.value, b.field))
	return b
}

// Numeric validates that the string contains only digits.
func (b *StrBuilder) Numeric() *StrBuilder {
	b.validations = append(b.validations, Numeric(b.value, b.field))
	return b
}

// AlphaUnicode validates that the string contains only Unicode letters.
func (b *StrBuilder) AlphaUnicode() *StrBuilder {
	b.validations = append(b.validations, AlphaUnicode(b.value, b.field))
	return b
}

// AlphaNumericUnicode validates that the string contains only Unicode letters and digits.
func (b *StrBuilder) AlphaNumericUnicode() *StrBuilder {
	b.validations = append(b.validations, AlphaNumericUnicode(b.value, b.field))
	return b
}

// ASCII validates that the string contains only ASCII characters.
func (b *StrBuilder) ASCII() *StrBuilder {
	b.validations = append(b.validations, ASCII(b.value, b.field))
	return b
}

// PrintableASCII validates that the string contains only printable ASCII.
func (b *StrBuilder) PrintableASCII() *StrBuilder {
	b.validations = append(b.validations, PrintableASCII(b.value, b.field))
	return b
}

// LowerCase validates that the string is entirely lowercase.
func (b *StrBuilder) LowerCase() *StrBuilder {
	b.validations = append(b.validations, LowerCase(b.value, b.field))
	return b
}

// UpperCase validates that the string is entirely uppercase.
func (b *StrBuilder) UpperCase() *StrBuilder {
	b.validations = append(b.validations, UpperCase(b.value, b.field))
	return b
}

// NoWhitespace validates that the string contains no whitespace.
func (b *StrBuilder) NoWhitespace() *StrBuilder {
	b.validations = append(b.validations, NoWhitespace(b.value, b.field))
	return b
}

// Trimmed validates that the string has no leading or trailing whitespace.
func (b *StrBuilder) Trimmed() *StrBuilder {
	b.validations = append(b.validations, Trimmed(b.value, b.field))
	return b
}

// SingleLine validates that the string contains no newlines.
func (b *StrBuilder) SingleLine() *StrBuilder {
	b.validations = append(b.validations, SingleLine(b.value, b.field))
	return b
}

// Identifier validates that the string is a valid identifier.
func (b *StrBuilder) Identifier() *StrBuilder {
	b.validations = append(b.validations, Identifier(b.value, b.field))
	return b
}

// Slug validates that the string is a valid URL slug.
func (b *StrBuilder) Slug() *StrBuilder {
	b.validations = append(b.validations, Slug(b.value, b.field))
	return b
}

// Email validates that the string is a valid email address.
func (b *StrBuilder) Email() *StrBuilder {
	b.validations = append(b.validations, Email(b.value, b.field))
	return b
}

// URL validates that the string is a valid URL.
func (b *StrBuilder) URL() *StrBuilder {
	b.validations = append(b.validations, URL(b.value, b.field))
	return b
}

// URLWithScheme validates that the string is a valid URL with one of the given schemes.
func (b *StrBuilder) URLWithScheme(schemes []string) *StrBuilder {
	b.validations = append(b.validations, URLWithScheme(b.value, schemes, b.field))
	return b
}

// HTTPOrHTTPS validates that the string is a valid HTTP or HTTPS URL.
func (b *StrBuilder) HTTPOrHTTPS() *StrBuilder {
	b.validations = append(b.validations, HTTPOrHTTPS(b.value, b.field))
	return b
}

// UUID validates that the string is a valid UUID.
func (b *StrBuilder) UUID() *StrBuilder {
	b.validations = append(b.validations, UUID(b.value, b.field))
	return b
}

// UUID4 validates that the string is a valid UUID v4.
func (b *StrBuilder) UUID4() *StrBuilder {
	b.validations = append(b.validations, UUID4(b.value, b.field))
	return b
}

// IP validates that the string is a valid IP address.
func (b *StrBuilder) IP() *StrBuilder {
	b.validations = append(b.validations, IP(b.value, b.field))
	return b
}

// IPv4 validates that the string is a valid IPv4 address.
func (b *StrBuilder) IPv4() *StrBuilder {
	b.validations = append(b.validations, IPv4(b.value, b.field))
	return b
}

// IPv6 validates that the string is a valid IPv6 address.
func (b *StrBuilder) IPv6() *StrBuilder {
	b.validations = append(b.validations, IPv6(b.value, b.field))
	return b
}

// CIDR validates that the string is a valid CIDR notation.
func (b *StrBuilder) CIDR() *StrBuilder {
	b.validations = append(b.validations, CIDR(b.value, b.field))
	return b
}

// MAC validates that the string is a valid MAC address.
func (b *StrBuilder) MAC() *StrBuilder {
	b.validations = append(b.validations, MAC(b.value, b.field))
	return b
}

// Hostname validates that the string is a valid hostname.
func (b *StrBuilder) Hostname() *StrBuilder {
	b.validations = append(b.validations, Hostname(b.value, b.field))
	return b
}

// Port validates that the string is a valid port number.
func (b *StrBuilder) Port() *StrBuilder {
	b.validations = append(b.validations, Port(b.value, b.field))
	return b
}

// HostPort validates that the string is a valid host:port combination.
func (b *StrBuilder) HostPort() *StrBuilder {
	b.validations = append(b.validations, HostPort(b.value, b.field))
	return b
}

// HexColor validates that the string is a valid hex color.
func (b *StrBuilder) HexColor() *StrBuilder {
	b.validations = append(b.validations, HexColor(b.value, b.field))
	return b
}

// HexColorFull validates that the string is a valid 6-digit hex color.
func (b *StrBuilder) HexColorFull() *StrBuilder {
	b.validations = append(b.validations, HexColorFull(b.value, b.field))
	return b
}

// Base64 validates that the string is valid base64.
func (b *StrBuilder) Base64() *StrBuilder {
	b.validations = append(b.validations, Base64(b.value, b.field))
	return b
}

// Base64URL validates that the string is valid URL-safe base64.
func (b *StrBuilder) Base64URL() *StrBuilder {
	b.validations = append(b.validations, Base64URL(b.value, b.field))
	return b
}

// JSON validates that the string is valid JSON.
func (b *StrBuilder) JSON() *StrBuilder {
	b.validations = append(b.validations, JSON(b.value, b.field))
	return b
}

// Semver validates that the string is a valid semantic version.
func (b *StrBuilder) Semver() *StrBuilder {
	b.validations = append(b.validations, Semver(b.value, b.field))
	return b
}

// E164 validates that the string is a valid E.164 phone number.
func (b *StrBuilder) E164() *StrBuilder {
	b.validations = append(b.validations, E164(b.value, b.field))
	return b
}

// CreditCard validates that the string is a valid credit card number.
func (b *StrBuilder) CreditCard() *StrBuilder {
	b.validations = append(b.validations, CreditCard(b.value, b.field))
	return b
}

// Latitude validates that the string is a valid latitude.
func (b *StrBuilder) Latitude() *StrBuilder {
	b.validations = append(b.validations, Latitude(b.value, b.field))
	return b
}

// Longitude validates that the string is a valid longitude.
func (b *StrBuilder) Longitude() *StrBuilder {
	b.validations = append(b.validations, Longitude(b.value, b.field))
	return b
}

// CountryCode2 validates that the string is a valid ISO 3166-1 alpha-2 country code.
func (b *StrBuilder) CountryCode2() *StrBuilder {
	b.validations = append(b.validations, CountryCode2(b.value, b.field))
	return b
}

// CountryCode3 validates that the string is a valid ISO 3166-1 alpha-3 country code.
func (b *StrBuilder) CountryCode3() *StrBuilder {
	b.validations = append(b.validations, CountryCode3(b.value, b.field))
	return b
}

// LanguageCode validates that the string is a valid ISO 639-1 language code.
func (b *StrBuilder) LanguageCode() *StrBuilder {
	b.validations = append(b.validations, LanguageCode(b.value, b.field))
	return b
}

// CurrencyCode validates that the string is a valid ISO 4217 currency code.
func (b *StrBuilder) CurrencyCode() *StrBuilder {
	b.validations = append(b.validations, CurrencyCode(b.value, b.field))
	return b
}

// Hex validates that the string contains only hexadecimal characters.
func (b *StrBuilder) Hex() *StrBuilder {
	b.validations = append(b.validations, Hex(b.value, b.field))
	return b
}

// DataURI validates that the string is a valid data URI.
func (b *StrBuilder) DataURI() *StrBuilder {
	b.validations = append(b.validations, DataURI(b.value, b.field))
	return b
}

// FilePath validates that the string is a valid file path.
func (b *StrBuilder) FilePath() *StrBuilder {
	b.validations = append(b.validations, FilePath(b.value, b.field))
	return b
}

// UnixPath validates that the string is a valid Unix path.
func (b *StrBuilder) UnixPath() *StrBuilder {
	b.validations = append(b.validations, UnixPath(b.value, b.field))
	return b
}

// -----------------------------------------------------------------------------
// Optional String Builder
// -----------------------------------------------------------------------------

// OptStrBuilder provides fluent validation for optional string pointers.
type OptStrBuilder struct {
	value       *string
	field       string
	validations []*Validation
	skip        bool
}

// OptStr creates a new optional string validation builder.
// If the pointer is nil, all validations are skipped (field is optional).
func OptStr(v *string, field string) *OptStrBuilder {
	return &OptStrBuilder{value: v, field: field, skip: v == nil}
}

// V returns the combined validation result.
// Returns nil if the value is nil (optional field not provided).
func (b *OptStrBuilder) V() *Validation {
	if b.skip {
		return nil
	}
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *OptStrBuilder) When(cond bool, fn func(*OptStrBuilder)) *OptStrBuilder {
	if cond && !b.skip {
		fn(b)
	}
	return b
}

// MinLen validates minimum string length.
func (b *OptStrBuilder) MinLen(n int) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, MinLen(*b.value, n, b.field))
	}
	return b
}

// MaxLen validates maximum string length.
func (b *OptStrBuilder) MaxLen(n int) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, MaxLen(*b.value, n, b.field))
	}
	return b
}

// Len validates exact string length.
func (b *OptStrBuilder) Len(n int) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Len(*b.value, n, b.field))
	}
	return b
}

// LenBetween validates string length is within a range.
func (b *OptStrBuilder) LenBetween(minLen, maxLen int) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, LenBetween(*b.value, minLen, maxLen, b.field))
	}
	return b
}

// Match validates that the string matches a pattern.
func (b *OptStrBuilder) Match(pattern *regexp.Regexp) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Match(*b.value, pattern, b.field))
	}
	return b
}

// NotMatch validates that the string does not match a pattern.
func (b *OptStrBuilder) NotMatch(pattern *regexp.Regexp) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, NotMatch(*b.value, pattern, b.field))
	}
	return b
}

// Prefix validates that the string starts with the given prefix.
func (b *OptStrBuilder) Prefix(prefix string) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Prefix(*b.value, prefix, b.field))
	}
	return b
}

// Suffix validates that the string ends with the given suffix.
func (b *OptStrBuilder) Suffix(suffix string) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Suffix(*b.value, suffix, b.field))
	}
	return b
}

// Contains validates that the string contains the substring.
func (b *OptStrBuilder) Contains(substr string) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Contains(*b.value, substr, b.field))
	}
	return b
}

// NotContains validates that the string does not contain the substring.
func (b *OptStrBuilder) NotContains(substr string) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, NotContains(*b.value, substr, b.field))
	}
	return b
}

// OneOf validates that the string is one of the allowed values.
func (b *OptStrBuilder) OneOf(allowed []string) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, OneOf(*b.value, allowed, b.field))
	}
	return b
}

// NotOneOf validates that the string is not one of the disallowed values.
func (b *OptStrBuilder) NotOneOf(disallowed []string) *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, NotOneOf(*b.value, disallowed, b.field))
	}
	return b
}

// Alpha validates that the string contains only ASCII letters.
func (b *OptStrBuilder) Alpha() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Alpha(*b.value, b.field))
	}
	return b
}

// AlphaNumeric validates that the string contains only ASCII letters and digits.
func (b *OptStrBuilder) AlphaNumeric() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, AlphaNumeric(*b.value, b.field))
	}
	return b
}

// Numeric validates that the string contains only digits.
func (b *OptStrBuilder) Numeric() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Numeric(*b.value, b.field))
	}
	return b
}

// LowerCase validates that the string is entirely lowercase.
func (b *OptStrBuilder) LowerCase() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, LowerCase(*b.value, b.field))
	}
	return b
}

// UpperCase validates that the string is entirely uppercase.
func (b *OptStrBuilder) UpperCase() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, UpperCase(*b.value, b.field))
	}
	return b
}

// Trimmed validates that the string has no leading or trailing whitespace.
func (b *OptStrBuilder) Trimmed() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Trimmed(*b.value, b.field))
	}
	return b
}

// SingleLine validates that the string contains no newlines.
func (b *OptStrBuilder) SingleLine() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, SingleLine(*b.value, b.field))
	}
	return b
}

// Slug validates that the string is a valid URL slug.
func (b *OptStrBuilder) Slug() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Slug(*b.value, b.field))
	}
	return b
}

// Email validates that the string is a valid email address.
func (b *OptStrBuilder) Email() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, Email(*b.value, b.field))
	}
	return b
}

// URL validates that the string is a valid URL.
func (b *OptStrBuilder) URL() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, URL(*b.value, b.field))
	}
	return b
}

// UUID validates that the string is a valid UUID.
func (b *OptStrBuilder) UUID() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, UUID(*b.value, b.field))
	}
	return b
}

// UUID4 validates that the string is a valid UUID v4.
func (b *OptStrBuilder) UUID4() *OptStrBuilder {
	if !b.skip {
		b.validations = append(b.validations, UUID4(*b.value, b.field))
	}
	return b
}

// -----------------------------------------------------------------------------
// Numeric Builder
// -----------------------------------------------------------------------------

// NumBuilder provides fluent validation for numeric values.
type NumBuilder[T constraints.Ordered] struct {
	value       T
	field       string
	validations []*Validation
}

// Num creates a new numeric validation builder.
func Num[T constraints.Ordered](v T, field string) *NumBuilder[T] {
	return &NumBuilder[T]{value: v, field: field}
}

// V returns the combined validation result.
func (b *NumBuilder[T]) V() *Validation {
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *NumBuilder[T]) When(cond bool, fn func(*NumBuilder[T])) *NumBuilder[T] {
	if cond {
		fn(b)
	}
	return b
}

// Min validates that the value is at least the minimum.
func (b *NumBuilder[T]) Min(minVal T) *NumBuilder[T] {
	b.validations = append(b.validations, Min(b.value, minVal, b.field))
	return b
}

// Max validates that the value is at most the maximum.
func (b *NumBuilder[T]) Max(maxVal T) *NumBuilder[T] {
	b.validations = append(b.validations, Max(b.value, maxVal, b.field))
	return b
}

// Between validates that the value is within a range (inclusive).
func (b *NumBuilder[T]) Between(minVal, maxVal T) *NumBuilder[T] {
	b.validations = append(b.validations, Between(b.value, minVal, maxVal, b.field))
	return b
}

// BetweenExclusive validates that the value is within a range (exclusive).
func (b *NumBuilder[T]) BetweenExclusive(minVal, maxVal T) *NumBuilder[T] {
	b.validations = append(b.validations, BetweenExclusive(b.value, minVal, maxVal, b.field))
	return b
}

// GreaterThan validates that the value is strictly greater than the threshold.
func (b *NumBuilder[T]) GreaterThan(threshold T) *NumBuilder[T] {
	b.validations = append(b.validations, GreaterThan(b.value, threshold, b.field))
	return b
}

// LessThan validates that the value is strictly less than the threshold.
func (b *NumBuilder[T]) LessThan(threshold T) *NumBuilder[T] {
	b.validations = append(b.validations, LessThan(b.value, threshold, b.field))
	return b
}

// GreaterThanOrEqual validates that the value is >= the threshold.
func (b *NumBuilder[T]) GreaterThanOrEqual(threshold T) *NumBuilder[T] {
	b.validations = append(b.validations, GreaterThanOrEqual(b.value, threshold, b.field))
	return b
}

// LessThanOrEqual validates that the value is <= the threshold.
func (b *NumBuilder[T]) LessThanOrEqual(threshold T) *NumBuilder[T] {
	b.validations = append(b.validations, LessThanOrEqual(b.value, threshold, b.field))
	return b
}

// OneOfValues validates that the value is one of the allowed values.
func (b *NumBuilder[T]) OneOfValues(allowed []T) *NumBuilder[T] {
	b.validations = append(b.validations, OneOfValues(b.value, allowed, b.field))
	return b
}

// NotOneOfValues validates that the value is not one of the disallowed values.
func (b *NumBuilder[T]) NotOneOfValues(disallowed []T) *NumBuilder[T] {
	b.validations = append(b.validations, NotOneOfValues(b.value, disallowed, b.field))
	return b
}

// -----------------------------------------------------------------------------
// Integer Builder (extends NumBuilder with integer-specific validators)
// -----------------------------------------------------------------------------

// IntBuilder provides fluent validation for integer values.
type IntBuilder[T Integer] struct {
	value       T
	field       string
	validations []*Validation
}

// Int creates a new integer validation builder.
func Int[T Integer](v T, field string) *IntBuilder[T] {
	return &IntBuilder[T]{value: v, field: field}
}

// V returns the combined validation result.
func (b *IntBuilder[T]) V() *Validation {
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *IntBuilder[T]) When(cond bool, fn func(*IntBuilder[T])) *IntBuilder[T] {
	if cond {
		fn(b)
	}
	return b
}

// Min validates that the value is at least the minimum.
func (b *IntBuilder[T]) Min(minVal T) *IntBuilder[T] {
	b.validations = append(b.validations, Min(b.value, minVal, b.field))
	return b
}

// Max validates that the value is at most the maximum.
func (b *IntBuilder[T]) Max(maxVal T) *IntBuilder[T] {
	b.validations = append(b.validations, Max(b.value, maxVal, b.field))
	return b
}

// Between validates that the value is within a range (inclusive).
func (b *IntBuilder[T]) Between(minVal, maxVal T) *IntBuilder[T] {
	b.validations = append(b.validations, Between(b.value, minVal, maxVal, b.field))
	return b
}

// Positive validates that the value is greater than zero.
func (b *IntBuilder[T]) Positive() *IntBuilder[T] {
	b.validations = append(b.validations, validation(func() error {
		if b.value <= 0 {
			return fieldErr(b.field, "must be positive")
		}
		return nil
	}(), b.field, "gt"))
	return b
}

// Negative validates that the value is less than zero.
func (b *IntBuilder[T]) Negative() *IntBuilder[T] {
	b.validations = append(b.validations, validation(func() error {
		if b.value >= 0 {
			return fieldErr(b.field, "must be negative")
		}
		return nil
	}(), b.field, "lt"))
	return b
}

// NonNegative validates that the value is zero or greater.
func (b *IntBuilder[T]) NonNegative() *IntBuilder[T] {
	b.validations = append(b.validations, validation(func() error {
		if b.value < 0 {
			return fieldErr(b.field, "must not be negative")
		}
		return nil
	}(), b.field, "gte"))
	return b
}

// NonPositive validates that the value is zero or less.
func (b *IntBuilder[T]) NonPositive() *IntBuilder[T] {
	b.validations = append(b.validations, validation(func() error {
		if b.value > 0 {
			return fieldErr(b.field, "must not be positive")
		}
		return nil
	}(), b.field, "lte"))
	return b
}

// Zero validates that the value is exactly zero.
func (b *IntBuilder[T]) Zero() *IntBuilder[T] {
	b.validations = append(b.validations, Zero(b.value, b.field))
	return b
}

// NonZero validates that the value is not zero.
func (b *IntBuilder[T]) NonZero() *IntBuilder[T] {
	b.validations = append(b.validations, NonZero(b.value, b.field))
	return b
}

// MultipleOf validates that the value is a multiple of the divisor.
func (b *IntBuilder[T]) MultipleOf(divisor T) *IntBuilder[T] {
	b.validations = append(b.validations, MultipleOf(b.value, divisor, b.field))
	return b
}

// Even validates that the value is even.
func (b *IntBuilder[T]) Even() *IntBuilder[T] {
	b.validations = append(b.validations, Even(b.value, b.field))
	return b
}

// Odd validates that the value is odd.
func (b *IntBuilder[T]) Odd() *IntBuilder[T] {
	b.validations = append(b.validations, Odd(b.value, b.field))
	return b
}

// -----------------------------------------------------------------------------
// Optional Numeric Builder
// -----------------------------------------------------------------------------

// OptNumBuilder provides fluent validation for optional numeric pointers.
type OptNumBuilder[T constraints.Ordered] struct {
	value       *T
	field       string
	validations []*Validation
	skip        bool
}

// OptNum creates a new optional numeric validation builder.
func OptNum[T constraints.Ordered](v *T, field string) *OptNumBuilder[T] {
	return &OptNumBuilder[T]{value: v, field: field, skip: v == nil}
}

// V returns the combined validation result.
func (b *OptNumBuilder[T]) V() *Validation {
	if b.skip {
		return nil
	}
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *OptNumBuilder[T]) When(cond bool, fn func(*OptNumBuilder[T])) *OptNumBuilder[T] {
	if cond && !b.skip {
		fn(b)
	}
	return b
}

// Min validates that the value is at least the minimum.
func (b *OptNumBuilder[T]) Min(minVal T) *OptNumBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, Min(*b.value, minVal, b.field))
	}
	return b
}

// Max validates that the value is at most the maximum.
func (b *OptNumBuilder[T]) Max(maxVal T) *OptNumBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, Max(*b.value, maxVal, b.field))
	}
	return b
}

// Between validates that the value is within a range (inclusive).
func (b *OptNumBuilder[T]) Between(minVal, maxVal T) *OptNumBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, Between(*b.value, minVal, maxVal, b.field))
	}
	return b
}

// GreaterThan validates that the value is strictly greater than the threshold.
func (b *OptNumBuilder[T]) GreaterThan(threshold T) *OptNumBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, GreaterThan(*b.value, threshold, b.field))
	}
	return b
}

// LessThan validates that the value is strictly less than the threshold.
func (b *OptNumBuilder[T]) LessThan(threshold T) *OptNumBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, LessThan(*b.value, threshold, b.field))
	}
	return b
}

// -----------------------------------------------------------------------------
// Slice Builders
// -----------------------------------------------------------------------------

// SliceBuilder provides fluent validation for slice values.
type SliceBuilder[T any] struct {
	value       []T
	field       string
	validations []*Validation
}

// Slice creates a new slice validation builder.
func Slice[T any](v []T, field string) *SliceBuilder[T] {
	return &SliceBuilder[T]{value: v, field: field}
}

// V returns the combined validation result.
func (b *SliceBuilder[T]) V() *Validation {
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *SliceBuilder[T]) When(cond bool, fn func(*SliceBuilder[T])) *SliceBuilder[T] {
	if cond {
		fn(b)
	}
	return b
}

// NotEmpty validates that the slice is not empty.
func (b *SliceBuilder[T]) NotEmpty() *SliceBuilder[T] {
	b.validations = append(b.validations, NotEmpty(b.value, b.field))
	return b
}

// Empty validates that the slice is empty.
func (b *SliceBuilder[T]) Empty() *SliceBuilder[T] {
	b.validations = append(b.validations, Empty(b.value, b.field))
	return b
}

// MinItems validates minimum slice length.
func (b *SliceBuilder[T]) MinItems(n int) *SliceBuilder[T] {
	b.validations = append(b.validations, MinItems(b.value, n, b.field))
	return b
}

// MaxItems validates maximum slice length.
func (b *SliceBuilder[T]) MaxItems(n int) *SliceBuilder[T] {
	b.validations = append(b.validations, MaxItems(b.value, n, b.field))
	return b
}

// ExactItems validates exact slice length.
func (b *SliceBuilder[T]) ExactItems(n int) *SliceBuilder[T] {
	b.validations = append(b.validations, ExactItems(b.value, n, b.field))
	return b
}

// ItemsBetween validates slice length is within a range.
func (b *SliceBuilder[T]) ItemsBetween(minItems, maxItems int) *SliceBuilder[T] {
	b.validations = append(b.validations, ItemsBetween(b.value, minItems, maxItems, b.field))
	return b
}

// Each applies a validation function to each element, collecting results.
// The function receives the element and auto-generated field name "field[i]".
// Returns *Validation for each element; non-nil results are collected.
func (b *SliceBuilder[T]) Each(fn func(v T, field string) *Validation) *SliceBuilder[T] {
	return b.EachV(fn)
}

// EachV applies a validation to each element, collecting results.
// The function receives the element and auto-generated field name.
func (b *SliceBuilder[T]) EachV(fn func(v T, field string) *Validation) *SliceBuilder[T] {
	for i, item := range b.value {
		elemField := fmt.Sprintf("%s[%d]", b.field, i)
		if v := fn(item, elemField); v != nil {
			b.validations = append(b.validations, v)
		}
	}
	return b
}

// -----------------------------------------------------------------------------
// String Slice Builder (with typed Each)
// -----------------------------------------------------------------------------

// StrSliceBuilder provides fluent validation for string slices.
type StrSliceBuilder struct {
	value       []string
	field       string
	validations []*Validation
}

// StrSlice creates a new string slice validation builder.
func StrSlice(v []string, field string) *StrSliceBuilder {
	return &StrSliceBuilder{value: v, field: field}
}

// V returns the combined validation result.
func (b *StrSliceBuilder) V() *Validation {
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *StrSliceBuilder) When(cond bool, fn func(*StrSliceBuilder)) *StrSliceBuilder {
	if cond {
		fn(b)
	}
	return b
}

// NotEmpty validates that the slice is not empty.
func (b *StrSliceBuilder) NotEmpty() *StrSliceBuilder {
	b.validations = append(b.validations, NotEmpty(b.value, b.field))
	return b
}

// MinItems validates minimum slice length.
func (b *StrSliceBuilder) MinItems(n int) *StrSliceBuilder {
	b.validations = append(b.validations, MinItems(b.value, n, b.field))
	return b
}

// MaxItems validates maximum slice length.
func (b *StrSliceBuilder) MaxItems(n int) *StrSliceBuilder {
	b.validations = append(b.validations, MaxItems(b.value, n, b.field))
	return b
}

// ItemsBetween validates slice length is within a range.
func (b *StrSliceBuilder) ItemsBetween(minItems, maxItems int) *StrSliceBuilder {
	b.validations = append(b.validations, ItemsBetween(b.value, minItems, maxItems, b.field))
	return b
}

// Unique validates that all elements are unique.
func (b *StrSliceBuilder) Unique() *StrSliceBuilder {
	b.validations = append(b.validations, Unique(b.value, b.field))
	return b
}

// Each applies validations to each element via a StrBuilder.
// Field names are auto-generated as "field[i]".
func (b *StrSliceBuilder) Each(fn func(*StrBuilder)) *StrSliceBuilder {
	for i, item := range b.value {
		elemField := fmt.Sprintf("%s[%d]", b.field, i)
		sb := &StrBuilder{value: item, field: elemField}
		fn(sb)
		if v := sb.V(); v != nil {
			b.validations = append(b.validations, v)
		}
	}
	return b
}

// AllMaxLen validates that all elements have at most n characters.
func (b *StrSliceBuilder) AllMaxLen(n int) *StrSliceBuilder {
	return b.Each(func(sb *StrBuilder) { sb.MaxLen(n) })
}

// AllMinLen validates that all elements have at least n characters.
func (b *StrSliceBuilder) AllMinLen(n int) *StrSliceBuilder {
	return b.Each(func(sb *StrBuilder) { sb.MinLen(n) })
}

// AllNotBlank validates that no element is blank.
func (b *StrSliceBuilder) AllNotBlank() *StrSliceBuilder {
	return b.Each(func(sb *StrBuilder) { sb.NotBlank() })
}

// -----------------------------------------------------------------------------
// Optional Integer Builder
// -----------------------------------------------------------------------------

// OptIntBuilder provides fluent validation for optional integer pointers.
type OptIntBuilder[T Integer] struct {
	value       *T
	field       string
	validations []*Validation
	skip        bool
}

// OptInt creates a new optional integer validation builder.
func OptInt[T Integer](v *T, field string) *OptIntBuilder[T] {
	return &OptIntBuilder[T]{value: v, field: field, skip: v == nil}
}

// V returns the combined validation result.
func (b *OptIntBuilder[T]) V() *Validation {
	if b.skip {
		return nil
	}
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *OptIntBuilder[T]) When(cond bool, fn func(*OptIntBuilder[T])) *OptIntBuilder[T] {
	if cond && !b.skip {
		fn(b)
	}
	return b
}

// Min validates that the value is at least the minimum.
func (b *OptIntBuilder[T]) Min(minVal T) *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, Min(*b.value, minVal, b.field))
	}
	return b
}

// Max validates that the value is at most the maximum.
func (b *OptIntBuilder[T]) Max(maxVal T) *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, Max(*b.value, maxVal, b.field))
	}
	return b
}

// Between validates that the value is within a range (inclusive).
func (b *OptIntBuilder[T]) Between(minVal, maxVal T) *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, Between(*b.value, minVal, maxVal, b.field))
	}
	return b
}

// Positive validates that the value is greater than zero.
func (b *OptIntBuilder[T]) Positive() *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, validation(func() error {
			if *b.value <= 0 {
				return fieldErr(b.field, "must be positive")
			}
			return nil
		}(), b.field, "gt"))
	}
	return b
}

// NonNegative validates that the value is zero or greater.
func (b *OptIntBuilder[T]) NonNegative() *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, validation(func() error {
			if *b.value < 0 {
				return fieldErr(b.field, "must not be negative")
			}
			return nil
		}(), b.field, "gte"))
	}
	return b
}

// NonZero validates that the value is not zero.
func (b *OptIntBuilder[T]) NonZero() *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, NonZero(*b.value, b.field))
	}
	return b
}

// MultipleOf validates that the value is a multiple of the divisor.
func (b *OptIntBuilder[T]) MultipleOf(divisor T) *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, MultipleOf(*b.value, divisor, b.field))
	}
	return b
}

// Even validates that the value is even.
func (b *OptIntBuilder[T]) Even() *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, Even(*b.value, b.field))
	}
	return b
}

// Odd validates that the value is odd.
func (b *OptIntBuilder[T]) Odd() *OptIntBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, Odd(*b.value, b.field))
	}
	return b
}

// -----------------------------------------------------------------------------
// Optional Slice Builder
// -----------------------------------------------------------------------------

// OptSliceBuilder provides fluent validation for optional slice pointers.
type OptSliceBuilder[T any] struct {
	value       *[]T
	field       string
	validations []*Validation
	skip        bool
}

// OptSlice creates a new optional slice validation builder.
func OptSlice[T any](v *[]T, field string) *OptSliceBuilder[T] {
	return &OptSliceBuilder[T]{value: v, field: field, skip: v == nil}
}

// V returns the combined validation result.
func (b *OptSliceBuilder[T]) V() *Validation {
	if b.skip {
		return nil
	}
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *OptSliceBuilder[T]) When(cond bool, fn func(*OptSliceBuilder[T])) *OptSliceBuilder[T] {
	if cond && !b.skip {
		fn(b)
	}
	return b
}

// NotEmpty validates that the slice is not empty.
func (b *OptSliceBuilder[T]) NotEmpty() *OptSliceBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, NotEmpty(*b.value, b.field))
	}
	return b
}

// MinItems validates minimum slice length.
func (b *OptSliceBuilder[T]) MinItems(n int) *OptSliceBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, MinItems(*b.value, n, b.field))
	}
	return b
}

// MaxItems validates maximum slice length.
func (b *OptSliceBuilder[T]) MaxItems(n int) *OptSliceBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, MaxItems(*b.value, n, b.field))
	}
	return b
}

// ItemsBetween validates slice length is within a range.
func (b *OptSliceBuilder[T]) ItemsBetween(minItems, maxItems int) *OptSliceBuilder[T] {
	if !b.skip {
		b.validations = append(b.validations, ItemsBetween(*b.value, minItems, maxItems, b.field))
	}
	return b
}

// EachV applies a validation to each element, collecting results.
func (b *OptSliceBuilder[T]) EachV(fn func(v T, field string) *Validation) *OptSliceBuilder[T] {
	if !b.skip {
		for i, item := range *b.value {
			elemField := fmt.Sprintf("%s[%d]", b.field, i)
			if v := fn(item, elemField); v != nil {
				b.validations = append(b.validations, v)
			}
		}
	}
	return b
}

// -----------------------------------------------------------------------------
// Optional String Slice Builder
// -----------------------------------------------------------------------------

// OptStrSliceBuilder provides fluent validation for optional string slice pointers.
type OptStrSliceBuilder struct {
	value       *[]string
	field       string
	validations []*Validation
	skip        bool
}

// OptStrSlice creates a new optional string slice validation builder.
func OptStrSlice(v *[]string, field string) *OptStrSliceBuilder {
	return &OptStrSliceBuilder{value: v, field: field, skip: v == nil}
}

// V returns the combined validation result.
func (b *OptStrSliceBuilder) V() *Validation {
	if b.skip {
		return nil
	}
	return combine(b.field, b.validations)
}

// When conditionally applies validations.
func (b *OptStrSliceBuilder) When(cond bool, fn func(*OptStrSliceBuilder)) *OptStrSliceBuilder {
	if cond && !b.skip {
		fn(b)
	}
	return b
}

// NotEmpty validates that the slice is not empty.
func (b *OptStrSliceBuilder) NotEmpty() *OptStrSliceBuilder {
	if !b.skip {
		b.validations = append(b.validations, NotEmpty(*b.value, b.field))
	}
	return b
}

// MinItems validates minimum slice length.
func (b *OptStrSliceBuilder) MinItems(n int) *OptStrSliceBuilder {
	if !b.skip {
		b.validations = append(b.validations, MinItems(*b.value, n, b.field))
	}
	return b
}

// MaxItems validates maximum slice length.
func (b *OptStrSliceBuilder) MaxItems(n int) *OptStrSliceBuilder {
	if !b.skip {
		b.validations = append(b.validations, MaxItems(*b.value, n, b.field))
	}
	return b
}

// Unique validates that all elements are unique.
func (b *OptStrSliceBuilder) Unique() *OptStrSliceBuilder {
	if !b.skip {
		b.validations = append(b.validations, Unique(*b.value, b.field))
	}
	return b
}

// Each applies validations to each element via a StrBuilder.
func (b *OptStrSliceBuilder) Each(fn func(*StrBuilder)) *OptStrSliceBuilder {
	if !b.skip {
		for i, item := range *b.value {
			elemField := fmt.Sprintf("%s[%d]", b.field, i)
			sb := &StrBuilder{value: item, field: elemField}
			fn(sb)
			if v := sb.V(); v != nil {
				b.validations = append(b.validations, v)
			}
		}
	}
	return b
}

// AllMaxLen validates that all elements have at most n characters.
func (b *OptStrSliceBuilder) AllMaxLen(n int) *OptStrSliceBuilder {
	return b.Each(func(sb *StrBuilder) { sb.MaxLen(n) })
}

// AllNotBlank validates that no element is blank.
func (b *OptStrSliceBuilder) AllNotBlank() *OptStrSliceBuilder {
	return b.Each(func(sb *StrBuilder) { sb.NotBlank() })
}
