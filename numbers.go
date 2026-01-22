package check

import (
	"golang.org/x/exp/constraints"
)

// Min validates that a value is at least the minimum.
func Min[T constraints.Ordered](v, minVal T, field string) error {
	if v < minVal {
		return fieldErrf(field, "must be at least %v", minVal)
	}
	return nil
}

// Max validates that a value is at most the maximum.
func Max[T constraints.Ordered](v, maxVal T, field string) error {
	if v > maxVal {
		return fieldErrf(field, "must be at most %v", maxVal)
	}
	return nil
}

// Between validates that a value is within a range (inclusive).
func Between[T constraints.Ordered](v, minVal, maxVal T, field string) error {
	if v < minVal || v > maxVal {
		return fieldErrf(field, "must be between %v and %v", minVal, maxVal)
	}
	return nil
}

// BetweenExclusive validates that a value is within a range (exclusive).
func BetweenExclusive[T constraints.Ordered](v, minVal, maxVal T, field string) error {
	if v <= minVal || v >= maxVal {
		return fieldErrf(field, "must be between %v and %v (exclusive)", minVal, maxVal)
	}
	return nil
}

// Signed is a constraint for signed numeric types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint for unsigned numeric types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float is a constraint for floating-point types.
type Float interface {
	~float32 | ~float64
}

// Integer is a constraint for all integer types.
type Integer interface {
	Signed | Unsigned
}

// Number is a constraint for all numeric types.
type Number interface {
	Integer | Float
}

// Positive validates that a value is greater than zero.
func Positive[T Signed | Float](v T, field string) error {
	if v <= 0 {
		return fieldErr(field, "must be positive")
	}
	return nil
}

// Negative validates that a value is less than zero.
func Negative[T Signed | Float](v T, field string) error {
	if v >= 0 {
		return fieldErr(field, "must be negative")
	}
	return nil
}

// NonNegative validates that a value is zero or greater.
func NonNegative[T Signed | Float](v T, field string) error {
	if v < 0 {
		return fieldErr(field, "must not be negative")
	}
	return nil
}

// NonPositive validates that a value is zero or less.
func NonPositive[T Signed | Float](v T, field string) error {
	if v > 0 {
		return fieldErr(field, "must not be positive")
	}
	return nil
}

// Zero validates that a value is exactly zero.
func Zero[T Number](v T, field string) error {
	if v != 0 {
		return fieldErr(field, "must be zero")
	}
	return nil
}

// NonZero validates that a value is not zero.
func NonZero[T Number](v T, field string) error {
	if v == 0 {
		return fieldErr(field, "must not be zero")
	}
	return nil
}

// MultipleOf validates that a value is a multiple of the given divisor.
func MultipleOf[T Integer](v, divisor T, field string) error {
	if divisor == 0 {
		return fieldErr(field, "divisor must not be zero")
	}
	if v%divisor != 0 {
		return fieldErrf(field, "must be a multiple of %v", divisor)
	}
	return nil
}

// Even validates that an integer value is even.
func Even[T Integer](v T, field string) error {
	if v%2 != 0 {
		return fieldErr(field, "must be even")
	}
	return nil
}

// Odd validates that an integer value is odd.
func Odd[T Integer](v T, field string) error {
	if v%2 == 0 {
		return fieldErr(field, "must be odd")
	}
	return nil
}

// OneOfValues validates that a value is one of the allowed values.
func OneOfValues[T comparable](v T, allowed []T, field string) error {
	for _, a := range allowed {
		if v == a {
			return nil
		}
	}
	return fieldErrf(field, "must be one of the allowed values")
}

// NotOneOfValues validates that a value is not one of the disallowed values.
func NotOneOfValues[T comparable](v T, disallowed []T, field string) error {
	for _, d := range disallowed {
		if v == d {
			return fieldErr(field, "must not be one of the disallowed values")
		}
	}
	return nil
}

// GreaterThan validates that a value is strictly greater than the threshold.
func GreaterThan[T constraints.Ordered](v, threshold T, field string) error {
	if v <= threshold {
		return fieldErrf(field, "must be greater than %v", threshold)
	}
	return nil
}

// LessThan validates that a value is strictly less than the threshold.
func LessThan[T constraints.Ordered](v, threshold T, field string) error {
	if v >= threshold {
		return fieldErrf(field, "must be less than %v", threshold)
	}
	return nil
}

// GreaterThanOrEqual validates that a value is greater than or equal to the threshold.
func GreaterThanOrEqual[T constraints.Ordered](v, threshold T, field string) error {
	if v < threshold {
		return fieldErrf(field, "must be greater than or equal to %v", threshold)
	}
	return nil
}

// LessThanOrEqual validates that a value is less than or equal to the threshold.
func LessThanOrEqual[T constraints.Ordered](v, threshold T, field string) error {
	if v > threshold {
		return fieldErrf(field, "must be less than or equal to %v", threshold)
	}
	return nil
}

// Percentage validates that a value is between 0 and 100.
func Percentage[T Number](v T, field string) error {
	if v < 0 || v > 100 {
		return fieldErr(field, "must be a percentage (0-100)")
	}
	return nil
}

// PortNumber validates that a value is a valid port number (1-65535).
func PortNumber(v int, field string) error {
	if v < 1 || v > 65535 {
		return fieldErr(field, "must be a valid port number (1-65535)")
	}
	return nil
}

// HTTPStatusCode validates that a value is a valid HTTP status code (100-599).
func HTTPStatusCode(v int, field string) error {
	if v < 100 || v > 599 {
		return fieldErr(field, "must be a valid HTTP status code (100-599)")
	}
	return nil
}
