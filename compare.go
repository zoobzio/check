package check

import (
	"golang.org/x/exp/constraints"
)

// Equal validates that two values are equal.
func Equal[T comparable](v, expected T, field string) *Validation {
	var err error
	if v != expected {
		err = fieldErrf(field, "must equal %v", expected)
	}
	return validation(err, field, "eq")
}

// NotEqual validates that two values are not equal.
func NotEqual[T comparable](v, other T, field string) *Validation {
	var err error
	if v == other {
		err = fieldErrf(field, "must not equal %v", other)
	}
	return validation(err, field, "ne")
}

// EqualField validates that a value equals another field's value.
// Useful for password confirmation, etc.
func EqualField[T comparable](v, other T, field, otherField string) *Validation {
	var err error
	if v != other {
		err = fieldErrf(field, "must equal %s", otherField)
	}
	return validation(err, field, "eqfield")
}

// NotEqualField validates that a value does not equal another field's value.
func NotEqualField[T comparable](v, other T, field, otherField string) *Validation {
	var err error
	if v == other {
		err = fieldErrf(field, "must not equal %s", otherField)
	}
	return validation(err, field, "nefield")
}

// GreaterThanField validates that a value is greater than another field's value.
func GreaterThanField[T constraints.Ordered](v, other T, field, otherField string) *Validation {
	var err error
	if v <= other {
		err = fieldErrf(field, "must be greater than %s", otherField)
	}
	return validation(err, field, "gtfield")
}

// LessThanField validates that a value is less than another field's value.
func LessThanField[T constraints.Ordered](v, other T, field, otherField string) *Validation {
	var err error
	if v >= other {
		err = fieldErrf(field, "must be less than %s", otherField)
	}
	return validation(err, field, "ltfield")
}

// GreaterThanOrEqualField validates that a value is greater than or equal to another field's value.
func GreaterThanOrEqualField[T constraints.Ordered](v, other T, field, otherField string) *Validation {
	var err error
	if v < other {
		err = fieldErrf(field, "must be greater than or equal to %s", otherField)
	}
	return validation(err, field, "gtefield")
}

// LessThanOrEqualField validates that a value is less than or equal to another field's value.
func LessThanOrEqualField[T constraints.Ordered](v, other T, field, otherField string) *Validation {
	var err error
	if v > other {
		err = fieldErrf(field, "must be less than or equal to %s", otherField)
	}
	return validation(err, field, "ltefield")
}
