# check

A zero-reflection validation library for Go. Explicit validation functions, no struct tags.

## Goals

- Simple function-based validation primitives
- No reflection, no struct tags
- Composable via `All()` for collecting errors
- Clear error messages with field names
- Minimal API surface

## API Design

```go
package check

// Core aggregator
func All(errs ...error) error

// String validators
func Required(v, field string) error
func MinLen(v string, min int, field string) error
func MaxLen(v string, max int, field string) error
func Len(v string, exact int, field string) error
func Match(v string, pattern *regexp.Regexp, field string) error

// Format validators
func Email(v, field string) error
func URL(v, field string) error
func UUID(v, field string) error

// Numeric validators
func Min[T constraints.Ordered](v, min T, field string) error
func Max[T constraints.Ordered](v, max T, field string) error
func Between[T constraints.Ordered](v, min, max T, field string) error
func Positive[T constraints.Signed | constraints.Float](v T, field string) error

// Slice validators
func NotEmpty[T any](v []T, field string) error
func MinItems[T any](v []T, min int, field string) error
func MaxItems[T any](v []T, max int, field string) error

// Pointer/optional validators
func NotNil[T any](v *T, field string) error
func NilOr[T any](v *T, fn func(T) error) error
```

## Usage

```go
func (r *RegisterRepositoryRequest) Validate() error {
    return check.All(
        check.Required(r.Owner, "owner"),
        check.MaxLen(r.Owner, 255, "owner"),
        check.Required(r.Name, "name"),
        check.MaxLen(r.Name, 255, "name"),
        check.Required(r.HTMLURL, "html_url"),
        check.URL(r.HTMLURL, "html_url"),
    )
}
```

## Error Format

```go
type FieldError struct {
    Field   string
    Message string
}

func (e *FieldError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

## File Structure

```
check/
├── check.go       # All() aggregator, FieldError type
├── strings.go     # String validators
├── formats.go     # Email, URL, UUID
├── numbers.go     # Numeric validators (generics)
├── slices.go      # Slice validators
└── pointers.go    # Pointer/optional validators
```

## Open Questions

1. Should `All()` return first error or collect all errors?
2. Should there be an `Errors` type that collects multiple `FieldError`s?
3. Should format validators (Email, URL) be strict or permissive?
