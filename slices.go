package check

// NotEmpty validates that a slice is not empty.
func NotEmpty[T any](v []T, field string) *Validation {
	var err error
	if len(v) == 0 {
		err = fieldErr(field, "must not be empty")
	}
	return validation(err, field, "required")
}

// Empty validates that a slice is empty.
func Empty[T any](v []T, field string) *Validation {
	var err error
	if len(v) != 0 {
		err = fieldErr(field, "must be empty")
	}
	return validation(err, field, "empty")
}

// MinItems validates minimum slice length.
func MinItems[T any](v []T, minCount int, field string) *Validation {
	var err error
	if len(v) < minCount {
		err = fieldErrf(field, "must have at least %d items", minCount)
	}
	return validation(err, field, "minitems")
}

// MaxItems validates maximum slice length.
func MaxItems[T any](v []T, maxCount int, field string) *Validation {
	var err error
	if len(v) > maxCount {
		err = fieldErrf(field, "must have at most %d items", maxCount)
	}
	return validation(err, field, "maxitems")
}

// ExactItems validates exact slice length.
func ExactItems[T any](v []T, count int, field string) *Validation {
	var err error
	if len(v) != count {
		err = fieldErrf(field, "must have exactly %d items", count)
	}
	return validation(err, field, "len")
}

// ItemsBetween validates slice length is within a range (inclusive).
func ItemsBetween[T any](v []T, minCount, maxCount int, field string) *Validation {
	var err error
	l := len(v)
	if l < minCount || l > maxCount {
		err = fieldErrf(field, "must have between %d and %d items", minCount, maxCount)
	}
	return validation(err, field, "minitems", "maxitems")
}

// Unique validates that all elements in a slice are unique.
func Unique[T comparable](v []T, field string) *Validation {
	var err error
	seen := make(map[T]struct{}, len(v))
	for _, item := range v {
		if _, exists := seen[item]; exists {
			err = fieldErr(field, "must have unique items")
			break
		}
		seen[item] = struct{}{}
	}
	return validation(err, field, "unique")
}

// SliceContains validates that a slice contains the given element.
func SliceContains[T comparable](v []T, elem T, field string) *Validation {
	var err error
	found := false
	for _, item := range v {
		if item == elem {
			found = true
			break
		}
	}
	if !found {
		err = fieldErrf(field, "must contain the required element")
	}
	return validation(err, field, "contains")
}

// SliceNotContains validates that a slice does not contain the given element.
func SliceNotContains[T comparable](v []T, elem T, field string) *Validation {
	var err error
	for _, item := range v {
		if item == elem {
			err = fieldErr(field, "must not contain the forbidden element")
			break
		}
	}
	return validation(err, field, "excludes")
}

// ContainsAll validates that a slice contains all the given elements.
func ContainsAll[T comparable](v []T, required []T, field string) *Validation {
	var err error
	set := make(map[T]struct{}, len(v))
	for _, item := range v {
		set[item] = struct{}{}
	}
	for _, req := range required {
		if _, exists := set[req]; !exists {
			err = fieldErr(field, "must contain all required elements")
			break
		}
	}
	return validation(err, field, "containsall")
}

// ContainsAny validates that a slice contains at least one of the given elements.
func ContainsAny[T comparable](v []T, options []T, field string) *Validation {
	var err error
	set := make(map[T]struct{}, len(v))
	for _, item := range v {
		set[item] = struct{}{}
	}
	found := false
	for _, opt := range options {
		if _, exists := set[opt]; exists {
			found = true
			break
		}
	}
	if !found {
		err = fieldErr(field, "must contain at least one of the required elements")
	}
	return validation(err, field, "containsany")
}

// ContainsNone validates that a slice contains none of the given elements.
func ContainsNone[T comparable](v []T, forbidden []T, field string) *Validation {
	var err error
	set := make(map[T]struct{}, len(v))
	for _, item := range v {
		set[item] = struct{}{}
	}
	for _, f := range forbidden {
		if _, exists := set[f]; exists {
			err = fieldErr(field, "must not contain any forbidden elements")
			break
		}
	}
	return validation(err, field, "excludesall")
}

// Each applies a validation function to each element in a slice.
// Returns a Result with all validations collected from each element.
func Each[T any](v []T, fn func(T, int) *Validation) *Result {
	validations := make([]*Validation, 0, len(v))
	for i, item := range v {
		val := fn(item, i)
		if val != nil {
			validations = append(validations, val)
		}
	}
	return All(validations...)
}

// EachValue applies a simple validation function (no index) to each element.
func EachValue[T any](v []T, fn func(T) *Validation) *Result {
	validations := make([]*Validation, 0, len(v))
	for _, item := range v {
		val := fn(item)
		if val != nil {
			validations = append(validations, val)
		}
	}
	return All(validations...)
}

// AllSatisfy validates that all elements satisfy a predicate.
func AllSatisfy[T any](v []T, pred func(T) bool, field, message string) *Validation {
	var err error
	for _, item := range v {
		if !pred(item) {
			err = fieldErr(field, message)
			break
		}
	}
	return validation(err, field, "all")
}

// AnySatisfies validates that at least one element satisfies a predicate.
func AnySatisfies[T any](v []T, pred func(T) bool, field, message string) *Validation {
	var err error
	found := false
	for _, item := range v {
		if pred(item) {
			found = true
			break
		}
	}
	if !found {
		err = fieldErr(field, message)
	}
	return validation(err, field, "any")
}

// NoneSatisfy validates that no elements satisfy a predicate.
func NoneSatisfy[T any](v []T, pred func(T) bool, field, message string) *Validation {
	var err error
	for _, item := range v {
		if pred(item) {
			err = fieldErr(field, message)
			break
		}
	}
	return validation(err, field, "none")
}

// Subset validates that all elements of v are in superset.
func Subset[T comparable](v, superset []T, field string) *Validation {
	var err error
	set := make(map[T]struct{}, len(superset))
	for _, item := range superset {
		set[item] = struct{}{}
	}
	for _, item := range v {
		if _, exists := set[item]; !exists {
			err = fieldErr(field, "must be a subset of the allowed values")
			break
		}
	}
	return validation(err, field, "subset")
}

// Disjoint validates that v shares no elements with other.
func Disjoint[T comparable](v, other []T, field string) *Validation {
	var err error
	set := make(map[T]struct{}, len(other))
	for _, item := range other {
		set[item] = struct{}{}
	}
	for _, item := range v {
		if _, exists := set[item]; exists {
			err = fieldErr(field, "must not share elements with the other set")
			break
		}
	}
	return validation(err, field, "disjoint")
}
