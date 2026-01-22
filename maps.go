package check

// NotEmptyMap validates that a map is not empty.
func NotEmptyMap[K comparable, V any](v map[K]V, field string) *Validation {
	var err error
	if len(v) == 0 {
		err = fieldErr(field, "must not be empty")
	}
	return validation(err, field, "required")
}

// EmptyMap validates that a map is empty.
func EmptyMap[K comparable, V any](v map[K]V, field string) *Validation {
	var err error
	if len(v) != 0 {
		err = fieldErr(field, "must be empty")
	}
	return validation(err, field, "empty")
}

// MinKeys validates minimum number of keys in a map.
func MinKeys[K comparable, V any](v map[K]V, minKeys int, field string) *Validation {
	var err error
	if len(v) < minKeys {
		err = fieldErrf(field, "must have at least %d keys", minKeys)
	}
	return validation(err, field, "minkeys")
}

// MaxKeys validates maximum number of keys in a map.
func MaxKeys[K comparable, V any](v map[K]V, maxKeys int, field string) *Validation {
	var err error
	if len(v) > maxKeys {
		err = fieldErrf(field, "must have at most %d keys", maxKeys)
	}
	return validation(err, field, "maxkeys")
}

// ExactKeys validates exact number of keys in a map.
func ExactKeys[K comparable, V any](v map[K]V, count int, field string) *Validation {
	var err error
	if len(v) != count {
		err = fieldErrf(field, "must have exactly %d keys", count)
	}
	return validation(err, field, "len")
}

// KeysBetween validates map size is within a range (inclusive).
func KeysBetween[K comparable, V any](v map[K]V, minKeys, maxKeys int, field string) *Validation {
	var err error
	l := len(v)
	if l < minKeys || l > maxKeys {
		err = fieldErrf(field, "must have between %d and %d keys", minKeys, maxKeys)
	}
	return validation(err, field, "minkeys", "maxkeys")
}

// HasKey validates that a map contains the given key.
func HasKey[K comparable, V any](v map[K]V, key K, field string) *Validation {
	var err error
	if _, exists := v[key]; !exists {
		err = fieldErrf(field, "must contain key %v", key)
	}
	return validation(err, field, "haskey")
}

// HasKeys validates that a map contains all the given keys.
func HasKeys[K comparable, V any](v map[K]V, keys []K, field string) *Validation {
	var err error
	for _, key := range keys {
		if _, exists := v[key]; !exists {
			err = fieldErr(field, "must contain all required keys")
			break
		}
	}
	return validation(err, field, "haskeys")
}

// HasAnyKey validates that a map contains at least one of the given keys.
func HasAnyKey[K comparable, V any](v map[K]V, keys []K, field string) *Validation {
	var err error
	found := false
	for _, key := range keys {
		if _, exists := v[key]; exists {
			found = true
			break
		}
	}
	if !found {
		err = fieldErr(field, "must contain at least one of the required keys")
	}
	return validation(err, field, "hasanykey")
}

// NotHasKey validates that a map does not contain the given key.
func NotHasKey[K comparable, V any](v map[K]V, key K, field string) *Validation {
	var err error
	if _, exists := v[key]; exists {
		err = fieldErrf(field, "must not contain key %v", key)
	}
	return validation(err, field, "nothaskey")
}

// NotHasKeys validates that a map does not contain any of the given keys.
func NotHasKeys[K comparable, V any](v map[K]V, keys []K, field string) *Validation {
	var err error
	for _, key := range keys {
		if _, exists := v[key]; exists {
			err = fieldErr(field, "must not contain any of the forbidden keys")
			break
		}
	}
	return validation(err, field, "nothaskeys")
}

// OnlyKeys validates that a map only contains keys from the allowed set.
func OnlyKeys[K comparable, V any](v map[K]V, allowed []K, field string) *Validation {
	var err error
	set := make(map[K]struct{}, len(allowed))
	for _, key := range allowed {
		set[key] = struct{}{}
	}
	for key := range v {
		if _, exists := set[key]; !exists {
			err = fieldErr(field, "must only contain allowed keys")
			break
		}
	}
	return validation(err, field, "onlykeys")
}

// EachKey applies a validation function to each key in a map.
func EachKey[K comparable, V any](v map[K]V, fn func(K) *Validation) *Result {
	validations := make([]*Validation, 0, len(v))
	for key := range v {
		val := fn(key)
		if val != nil {
			validations = append(validations, val)
		}
	}
	return All(validations...)
}

// EachMapValue applies a validation function to each value in a map.
func EachMapValue[K comparable, V any](v map[K]V, fn func(V) *Validation) *Result {
	validations := make([]*Validation, 0, len(v))
	for _, val := range v {
		result := fn(val)
		if result != nil {
			validations = append(validations, result)
		}
	}
	return All(validations...)
}

// EachEntry applies a validation function to each key-value pair in a map.
func EachEntry[K comparable, V any](v map[K]V, fn func(K, V) *Validation) *Result {
	validations := make([]*Validation, 0, len(v))
	for key, val := range v {
		result := fn(key, val)
		if result != nil {
			validations = append(validations, result)
		}
	}
	return All(validations...)
}

// UniqueValues validates that all values in a map are unique.
func UniqueValues[K, V comparable](v map[K]V, field string) *Validation {
	var err error
	seen := make(map[V]struct{}, len(v))
	for _, val := range v {
		if _, exists := seen[val]; exists {
			err = fieldErr(field, "must have unique values")
			break
		}
		seen[val] = struct{}{}
	}
	return validation(err, field, "unique")
}
