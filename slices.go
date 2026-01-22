package check

// NotEmpty validates that a slice is not empty.
func NotEmpty[T any](v []T, field string) error {
	if len(v) == 0 {
		return fieldErr(field, "must not be empty")
	}
	return nil
}

// Empty validates that a slice is empty.
func Empty[T any](v []T, field string) error {
	if len(v) != 0 {
		return fieldErr(field, "must be empty")
	}
	return nil
}

// MinItems validates minimum slice length.
func MinItems[T any](v []T, minCount int, field string) error {
	if len(v) < minCount {
		return fieldErrf(field, "must have at least %d items", minCount)
	}
	return nil
}

// MaxItems validates maximum slice length.
func MaxItems[T any](v []T, maxCount int, field string) error {
	if len(v) > maxCount {
		return fieldErrf(field, "must have at most %d items", maxCount)
	}
	return nil
}

// ExactItems validates exact slice length.
func ExactItems[T any](v []T, count int, field string) error {
	if len(v) != count {
		return fieldErrf(field, "must have exactly %d items", count)
	}
	return nil
}

// ItemsBetween validates slice length is within a range (inclusive).
func ItemsBetween[T any](v []T, minCount, maxCount int, field string) error {
	l := len(v)
	if l < minCount || l > maxCount {
		return fieldErrf(field, "must have between %d and %d items", minCount, maxCount)
	}
	return nil
}

// Unique validates that all elements in a slice are unique.
func Unique[T comparable](v []T, field string) error {
	seen := make(map[T]struct{}, len(v))
	for _, item := range v {
		if _, exists := seen[item]; exists {
			return fieldErr(field, "must have unique items")
		}
		seen[item] = struct{}{}
	}
	return nil
}

// SliceContains validates that a slice contains the given element.
func SliceContains[T comparable](v []T, elem T, field string) error {
	for _, item := range v {
		if item == elem {
			return nil
		}
	}
	return fieldErrf(field, "must contain the required element")
}

// SliceNotContains validates that a slice does not contain the given element.
func SliceNotContains[T comparable](v []T, elem T, field string) error {
	for _, item := range v {
		if item == elem {
			return fieldErr(field, "must not contain the forbidden element")
		}
	}
	return nil
}

// ContainsAll validates that a slice contains all the given elements.
func ContainsAll[T comparable](v []T, required []T, field string) error {
	set := make(map[T]struct{}, len(v))
	for _, item := range v {
		set[item] = struct{}{}
	}
	for _, req := range required {
		if _, exists := set[req]; !exists {
			return fieldErr(field, "must contain all required elements")
		}
	}
	return nil
}

// ContainsAny validates that a slice contains at least one of the given elements.
func ContainsAny[T comparable](v []T, options []T, field string) error {
	set := make(map[T]struct{}, len(v))
	for _, item := range v {
		set[item] = struct{}{}
	}
	for _, opt := range options {
		if _, exists := set[opt]; exists {
			return nil
		}
	}
	return fieldErr(field, "must contain at least one of the required elements")
}

// ContainsNone validates that a slice contains none of the given elements.
func ContainsNone[T comparable](v []T, forbidden []T, field string) error {
	set := make(map[T]struct{}, len(v))
	for _, item := range v {
		set[item] = struct{}{}
	}
	for _, f := range forbidden {
		if _, exists := set[f]; exists {
			return fieldErr(field, "must not contain any forbidden elements")
		}
	}
	return nil
}

// Each applies a validation function to each element in a slice.
// Returns all errors collected from validating each element.
func Each[T any](v []T, fn func(T, int) error) error {
	var errs Errors
	for i, item := range v {
		if err := fn(item, i); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// EachValue applies a simple validation function (no index) to each element.
func EachValue[T any](v []T, fn func(T) error) error {
	var errs Errors
	for _, item := range v {
		if err := fn(item); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// AllSatisfy validates that all elements satisfy a predicate.
func AllSatisfy[T any](v []T, pred func(T) bool, field, message string) error {
	for _, item := range v {
		if !pred(item) {
			return fieldErr(field, message)
		}
	}
	return nil
}

// AnySatisfies validates that at least one element satisfies a predicate.
func AnySatisfies[T any](v []T, pred func(T) bool, field, message string) error {
	for _, item := range v {
		if pred(item) {
			return nil
		}
	}
	return fieldErr(field, message)
}

// NoneSatisfy validates that no elements satisfy a predicate.
func NoneSatisfy[T any](v []T, pred func(T) bool, field, message string) error {
	for _, item := range v {
		if pred(item) {
			return fieldErr(field, message)
		}
	}
	return nil
}

// Subset validates that all elements of v are in superset.
func Subset[T comparable](v, superset []T, field string) error {
	set := make(map[T]struct{}, len(superset))
	for _, item := range superset {
		set[item] = struct{}{}
	}
	for _, item := range v {
		if _, exists := set[item]; !exists {
			return fieldErr(field, "must be a subset of the allowed values")
		}
	}
	return nil
}

// Disjoint validates that v shares no elements with other.
func Disjoint[T comparable](v, other []T, field string) error {
	set := make(map[T]struct{}, len(other))
	for _, item := range other {
		set[item] = struct{}{}
	}
	for _, item := range v {
		if _, exists := set[item]; exists {
			return fieldErr(field, "must not share elements with the other set")
		}
	}
	return nil
}
