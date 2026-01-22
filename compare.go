package check

import (
	"golang.org/x/exp/constraints"
)

// Equal validates that two values are equal.
func Equal[T comparable](v, expected T, field string) error {
	if v != expected {
		return fieldErrf(field, "must equal %v", expected)
	}
	return nil
}

// NotEqual validates that two values are not equal.
func NotEqual[T comparable](v, other T, field string) error {
	if v == other {
		return fieldErrf(field, "must not equal %v", other)
	}
	return nil
}

// EqualField validates that a value equals another field's value.
// Useful for password confirmation, etc.
func EqualField[T comparable](v, other T, field, otherField string) error {
	if v != other {
		return fieldErrf(field, "must equal %s", otherField)
	}
	return nil
}

// NotEqualField validates that a value does not equal another field's value.
func NotEqualField[T comparable](v, other T, field, otherField string) error {
	if v == other {
		return fieldErrf(field, "must not equal %s", otherField)
	}
	return nil
}

// GreaterThanField validates that a value is greater than another field's value.
func GreaterThanField[T constraints.Ordered](v, other T, field, otherField string) error {
	if v <= other {
		return fieldErrf(field, "must be greater than %s", otherField)
	}
	return nil
}

// LessThanField validates that a value is less than another field's value.
func LessThanField[T constraints.Ordered](v, other T, field, otherField string) error {
	if v >= other {
		return fieldErrf(field, "must be less than %s", otherField)
	}
	return nil
}

// GreaterThanOrEqualField validates that a value is greater than or equal to another field's value.
func GreaterThanOrEqualField[T constraints.Ordered](v, other T, field, otherField string) error {
	if v < other {
		return fieldErrf(field, "must be greater than or equal to %s", otherField)
	}
	return nil
}

// LessThanOrEqualField validates that a value is less than or equal to another field's value.
func LessThanOrEqualField[T constraints.Ordered](v, other T, field, otherField string) error {
	if v > other {
		return fieldErrf(field, "must be less than or equal to %s", otherField)
	}
	return nil
}
