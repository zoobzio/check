package check

// NotNil validates that a pointer is not nil.
func NotNil[T any](v *T, field string) *Validation {
	var err error
	if v == nil {
		err = fieldErr(field, "must not be nil")
	}
	return validation(err, field, "required")
}

// Nil validates that a pointer is nil.
func Nil[T any](v *T, field string) *Validation {
	var err error
	if v != nil {
		err = fieldErr(field, "must be nil")
	}
	return validation(err, field, "nil")
}

// NilOr validates a pointer value if it's not nil.
// If the pointer is nil, returns nil (no validation applied - field is optional).
// If the pointer is not nil, applies the given validation function.
func NilOr[T any](v *T, fn func(T) *Validation) *Validation {
	if v == nil {
		return nil
	}
	return fn(*v)
}

// NilOrField is like NilOr but includes field context in the error.
func NilOrField[T any](v *T, fn func(T, string) *Validation, field string) *Validation {
	if v == nil {
		return nil
	}
	return fn(*v, field)
}

// RequiredPtr validates that a pointer is not nil and applies validation to its value.
// Reports "required" validator, plus any validators from the inner function.
func RequiredPtr[T any](v *T, fn func(T) *Validation, field string) *Validation {
	if v == nil {
		return validation(fieldErr(field, "is required"), field, "required")
	}

	inner := fn(*v)
	if inner == nil {
		return validation(nil, field, "required")
	}

	// Combine required with inner validators
	validators := append([]string{"required"}, inner.validators...)
	return validation(inner.err, field, validators...)
}

// RequiredPtrField validates that a pointer is not nil and applies a field-aware validation.
func RequiredPtrField[T any](v *T, fn func(T, string) *Validation, field string) *Validation {
	if v == nil {
		return validation(fieldErr(field, "is required"), field, "required")
	}

	inner := fn(*v, field)
	if inner == nil {
		return validation(nil, field, "required")
	}

	// Combine required with inner validators
	validators := append([]string{"required"}, inner.validators...)
	return validation(inner.err, field, validators...)
}

// DefaultOr uses a default value if the pointer is nil, then validates.
func DefaultOr[T any](v *T, defaultVal T, fn func(T) *Validation) *Validation {
	val := defaultVal
	if v != nil {
		val = *v
	}
	return fn(val)
}

// Deref safely dereferences a pointer, returning the zero value if nil.
func Deref[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// DerefOr safely dereferences a pointer, returning the default if nil.
func DerefOr[T any](v *T, defaultVal T) T {
	if v == nil {
		return defaultVal
	}
	return *v
}

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// NotNilInterface validates that an interface value is not nil.
func NotNilInterface(v any, field string) *Validation {
	var err error
	if v == nil {
		err = fieldErr(field, "must not be nil")
	}
	return validation(err, field, "required")
}
