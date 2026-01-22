package check

import (
	"golang.org/x/exp/constraints"
)

// Min validates that a value is at least the minimum.
func Min[T constraints.Ordered](v, minVal T, field string) *Validation {
	var err error
	if v < minVal {
		err = fieldErrf(field, "must be at least %v", minVal)
	}
	return validation(err, field, "min")
}

// Max validates that a value is at most the maximum.
func Max[T constraints.Ordered](v, maxVal T, field string) *Validation {
	var err error
	if v > maxVal {
		err = fieldErrf(field, "must be at most %v", maxVal)
	}
	return validation(err, field, "max")
}

// Between validates that a value is within a range (inclusive).
func Between[T constraints.Ordered](v, minVal, maxVal T, field string) *Validation {
	var err error
	if v < minVal || v > maxVal {
		err = fieldErrf(field, "must be between %v and %v", minVal, maxVal)
	}
	return validation(err, field, "min", "max")
}

// BetweenExclusive validates that a value is within a range (exclusive).
func BetweenExclusive[T constraints.Ordered](v, minVal, maxVal T, field string) *Validation {
	var err error
	if v <= minVal || v >= maxVal {
		err = fieldErrf(field, "must be between %v and %v (exclusive)", minVal, maxVal)
	}
	return validation(err, field, "gt", "lt")
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
func Positive[T Signed | Float](v T, field string) *Validation {
	var err error
	if v <= 0 {
		err = fieldErr(field, "must be positive")
	}
	return validation(err, field, "gt")
}

// Negative validates that a value is less than zero.
func Negative[T Signed | Float](v T, field string) *Validation {
	var err error
	if v >= 0 {
		err = fieldErr(field, "must be negative")
	}
	return validation(err, field, "lt")
}

// NonNegative validates that a value is zero or greater.
func NonNegative[T Signed | Float](v T, field string) *Validation {
	var err error
	if v < 0 {
		err = fieldErr(field, "must not be negative")
	}
	return validation(err, field, "gte")
}

// NonPositive validates that a value is zero or less.
func NonPositive[T Signed | Float](v T, field string) *Validation {
	var err error
	if v > 0 {
		err = fieldErr(field, "must not be positive")
	}
	return validation(err, field, "lte")
}

// Zero validates that a value is exactly zero.
func Zero[T Number](v T, field string) *Validation {
	var err error
	if v != 0 {
		err = fieldErr(field, "must be zero")
	}
	return validation(err, field, "eq")
}

// NonZero validates that a value is not zero.
func NonZero[T Number](v T, field string) *Validation {
	var err error
	if v == 0 {
		err = fieldErr(field, "must not be zero")
	}
	return validation(err, field, "ne")
}

// MultipleOf validates that a value is a multiple of the given divisor.
func MultipleOf[T Integer](v, divisor T, field string) *Validation {
	var err error
	if divisor == 0 {
		err = fieldErr(field, "divisor must not be zero")
	} else if v%divisor != 0 {
		err = fieldErrf(field, "must be a multiple of %v", divisor)
	}
	return validation(err, field, "multipleof")
}

// Even validates that an integer value is even.
func Even[T Integer](v T, field string) *Validation {
	var err error
	if v%2 != 0 {
		err = fieldErr(field, "must be even")
	}
	return validation(err, field, "even")
}

// Odd validates that an integer value is odd.
func Odd[T Integer](v T, field string) *Validation {
	var err error
	if v%2 == 0 {
		err = fieldErr(field, "must be odd")
	}
	return validation(err, field, "odd")
}

// OneOfValues validates that a value is one of the allowed values.
func OneOfValues[T comparable](v T, allowed []T, field string) *Validation {
	var err error
	found := false
	for _, a := range allowed {
		if v == a {
			found = true
			break
		}
	}
	if !found {
		err = fieldErrf(field, "must be one of the allowed values")
	}
	return validation(err, field, "oneof")
}

// NotOneOfValues validates that a value is not one of the disallowed values.
func NotOneOfValues[T comparable](v T, disallowed []T, field string) *Validation {
	var err error
	for _, d := range disallowed {
		if v == d {
			err = fieldErr(field, "must not be one of the disallowed values")
			break
		}
	}
	return validation(err, field, "notoneof")
}

// GreaterThan validates that a value is strictly greater than the threshold.
func GreaterThan[T constraints.Ordered](v, threshold T, field string) *Validation {
	var err error
	if v <= threshold {
		err = fieldErrf(field, "must be greater than %v", threshold)
	}
	return validation(err, field, "gt")
}

// LessThan validates that a value is strictly less than the threshold.
func LessThan[T constraints.Ordered](v, threshold T, field string) *Validation {
	var err error
	if v >= threshold {
		err = fieldErrf(field, "must be less than %v", threshold)
	}
	return validation(err, field, "lt")
}

// GreaterThanOrEqual validates that a value is greater than or equal to the threshold.
func GreaterThanOrEqual[T constraints.Ordered](v, threshold T, field string) *Validation {
	var err error
	if v < threshold {
		err = fieldErrf(field, "must be greater than or equal to %v", threshold)
	}
	return validation(err, field, "gte")
}

// LessThanOrEqual validates that a value is less than or equal to the threshold.
func LessThanOrEqual[T constraints.Ordered](v, threshold T, field string) *Validation {
	var err error
	if v > threshold {
		err = fieldErrf(field, "must be less than or equal to %v", threshold)
	}
	return validation(err, field, "lte")
}

// Percentage validates that a value is between 0 and 100.
func Percentage[T Number](v T, field string) *Validation {
	var err error
	if v < 0 || v > 100 {
		err = fieldErr(field, "must be a percentage (0-100)")
	}
	return validation(err, field, "min", "max")
}

// PortNumber validates that a value is a valid port number (1-65535).
func PortNumber(v int, field string) *Validation {
	var err error
	if v < 1 || v > 65535 {
		err = fieldErr(field, "must be a valid port number (1-65535)")
	}
	return validation(err, field, "port")
}

// HTTPStatusCode validates that a value is a valid HTTP status code (100-599).
func HTTPStatusCode(v int, field string) *Validation {
	var err error
	if v < 100 || v > 599 {
		err = fieldErr(field, "must be a valid HTTP status code (100-599)")
	}
	return validation(err, field, "httpstatus")
}
