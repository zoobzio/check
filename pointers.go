package check

// NotNil validates that a pointer is not nil.
func NotNil[T any](v *T, field string) error {
	if v == nil {
		return fieldErr(field, "must not be nil")
	}
	return nil
}

// Nil validates that a pointer is nil.
func Nil[T any](v *T, field string) error {
	if v != nil {
		return fieldErr(field, "must be nil")
	}
	return nil
}

// NilOr validates a pointer value if it's not nil.
// If the pointer is nil, validation passes (the field is optional).
// If the pointer is not nil, applies the given validation function.
func NilOr[T any](v *T, fn func(T) error) error {
	if v == nil {
		return nil
	}
	return fn(*v)
}

// NilOrField is like NilOr but includes field context in the error.
func NilOrField[T any](v *T, fn func(T, string) error, field string) error {
	if v == nil {
		return nil
	}
	return fn(*v, field)
}

// RequiredPtr validates that a pointer is not nil and applies validation to its value.
func RequiredPtr[T any](v *T, fn func(T) error, field string) error {
	if v == nil {
		return fieldErr(field, "is required")
	}
	return fn(*v)
}

// RequiredPtrField validates that a pointer is not nil and applies a field-aware validation.
func RequiredPtrField[T any](v *T, fn func(T, string) error, field string) error {
	if v == nil {
		return fieldErr(field, "is required")
	}
	return fn(*v, field)
}

// DefaultOr uses a default value if the pointer is nil, then validates.
func DefaultOr[T any](v *T, defaultVal T, fn func(T) error) error {
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
// Note: This checks both typed nil and untyped nil.
func NotNilInterface(v any, field string) error {
	if v == nil {
		return fieldErr(field, "must not be nil")
	}
	return nil
}
