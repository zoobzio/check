package check

// NotEmptyMap validates that a map is not empty.
func NotEmptyMap[K comparable, V any](v map[K]V, field string) error {
	if len(v) == 0 {
		return fieldErr(field, "must not be empty")
	}
	return nil
}

// EmptyMap validates that a map is empty.
func EmptyMap[K comparable, V any](v map[K]V, field string) error {
	if len(v) != 0 {
		return fieldErr(field, "must be empty")
	}
	return nil
}

// MinKeys validates minimum number of keys in a map.
func MinKeys[K comparable, V any](v map[K]V, minKeys int, field string) error {
	if len(v) < minKeys {
		return fieldErrf(field, "must have at least %d keys", minKeys)
	}
	return nil
}

// MaxKeys validates maximum number of keys in a map.
func MaxKeys[K comparable, V any](v map[K]V, maxKeys int, field string) error {
	if len(v) > maxKeys {
		return fieldErrf(field, "must have at most %d keys", maxKeys)
	}
	return nil
}

// ExactKeys validates exact number of keys in a map.
func ExactKeys[K comparable, V any](v map[K]V, count int, field string) error {
	if len(v) != count {
		return fieldErrf(field, "must have exactly %d keys", count)
	}
	return nil
}

// KeysBetween validates map size is within a range (inclusive).
func KeysBetween[K comparable, V any](v map[K]V, minKeys, maxKeys int, field string) error {
	l := len(v)
	if l < minKeys || l > maxKeys {
		return fieldErrf(field, "must have between %d and %d keys", minKeys, maxKeys)
	}
	return nil
}

// HasKey validates that a map contains the given key.
func HasKey[K comparable, V any](v map[K]V, key K, field string) error {
	if _, exists := v[key]; !exists {
		return fieldErrf(field, "must contain key %v", key)
	}
	return nil
}

// HasKeys validates that a map contains all the given keys.
func HasKeys[K comparable, V any](v map[K]V, keys []K, field string) error {
	for _, key := range keys {
		if _, exists := v[key]; !exists {
			return fieldErr(field, "must contain all required keys")
		}
	}
	return nil
}

// HasAnyKey validates that a map contains at least one of the given keys.
func HasAnyKey[K comparable, V any](v map[K]V, keys []K, field string) error {
	for _, key := range keys {
		if _, exists := v[key]; exists {
			return nil
		}
	}
	return fieldErr(field, "must contain at least one of the required keys")
}

// NotHasKey validates that a map does not contain the given key.
func NotHasKey[K comparable, V any](v map[K]V, key K, field string) error {
	if _, exists := v[key]; exists {
		return fieldErrf(field, "must not contain key %v", key)
	}
	return nil
}

// NotHasKeys validates that a map does not contain any of the given keys.
func NotHasKeys[K comparable, V any](v map[K]V, keys []K, field string) error {
	for _, key := range keys {
		if _, exists := v[key]; exists {
			return fieldErr(field, "must not contain any of the forbidden keys")
		}
	}
	return nil
}

// OnlyKeys validates that a map only contains keys from the allowed set.
func OnlyKeys[K comparable, V any](v map[K]V, allowed []K, field string) error {
	set := make(map[K]struct{}, len(allowed))
	for _, key := range allowed {
		set[key] = struct{}{}
	}
	for key := range v {
		if _, exists := set[key]; !exists {
			return fieldErr(field, "must only contain allowed keys")
		}
	}
	return nil
}

// EachKey applies a validation function to each key in a map.
func EachKey[K comparable, V any](v map[K]V, fn func(K) error) error {
	var errs Errors
	for key := range v {
		if err := fn(key); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// EachMapValue applies a validation function to each value in a map.
func EachMapValue[K comparable, V any](v map[K]V, fn func(V) error) error {
	var errs Errors
	for _, val := range v {
		if err := fn(val); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// EachEntry applies a validation function to each key-value pair in a map.
func EachEntry[K comparable, V any](v map[K]V, fn func(K, V) error) error {
	var errs Errors
	for key, val := range v {
		if err := fn(key, val); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// UniqueValues validates that all values in a map are unique.
func UniqueValues[K, V comparable](v map[K]V, field string) error {
	seen := make(map[V]struct{}, len(v))
	for _, val := range v {
		if _, exists := seen[val]; exists {
			return fieldErr(field, "must have unique values")
		}
		seen[val] = struct{}{}
	}
	return nil
}
