# check

[![CI Status](https://github.com/zoobzio/check/workflows/CI/badge.svg)](https://github.com/zoobzio/check/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/zoobzio/check/graph/badge.svg?branch=main)](https://codecov.io/gh/zoobzio/check)
[![Go Report Card](https://goreportcard.com/badge/github.com/zoobzio/check)](https://goreportcard.com/report/github.com/zoobzio/check)
[![CodeQL](https://github.com/zoobzio/check/workflows/CodeQL/badge.svg)](https://github.com/zoobzio/check/security/code-scanning)
[![Go Reference](https://pkg.go.dev/badge/github.com/zoobzio/check.svg)](https://pkg.go.dev/github.com/zoobzio/check)
[![License](https://img.shields.io/github/license/zoobzio/check)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod-go-version/zoobzio/check)](go.mod)
[![Release](https://img.shields.io/github/v/release/zoobzio/check)](https://github.com/zoobzio/check/releases)

Zero-reflection validation primitives for Go.

Explicit validation functions, no struct tags, fully composable.

## Just Functions

```go
type User struct {
    Email    string
    Age      int
    Username string
}

func (u User) Validate() *check.Result {
    return check.All(
        check.Required(u.Email, "email"),
        check.Email(u.Email, "email"),
        check.Between(u.Age, 13, 120, "age"),
        check.LenBetween(u.Username, 3, 20, "username"),
        check.Slug(u.Username, "username"),
    )
}
```

No magic, no reflection, and no runtime struct-tag parsing—just functions that return validation results.

```go
r := user.Validate()
if r.Err() != nil {
    fmt.Println(r.Err())
    // email: must be a valid email address; age: must be between 13 and 120

    for _, fe := range check.GetFieldErrors(r) {
        fmt.Printf("%s: %s\n", fe.Field, fe.Message)
    }
    // email: must be a valid email address
    // age: must be between 13 and 120
}
```

Validation logic lives where you can see it, test it, and refactor it.

## Install

```bash
go get github.com/zoobzio/check
```

Requires Go 1.24+.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/zoobzio/check"
)

type Order struct {
    ID       string
    Total    float64
    Quantity int
    Email    string
}

func (o Order) Validate() *check.Result {
    return check.All(
        check.Required(o.ID, "id"),
        check.UUID(o.ID, "id"),
        check.Positive(o.Total, "total"),
        check.Between(o.Quantity, 1, 100, "quantity"),
        check.Email(o.Email, "email"),
    )
}

func main() {
    order := Order{
        ID:       "not-a-uuid",
        Total:    -50.00,
        Quantity: 0,
        Email:    "invalid",
    }

    r := order.Validate()
    if r.Err() != nil {
        // All errors collected
        fmt.Println(r.Err())
        // id: must be a valid UUID; total: must be positive; quantity: must be between 1 and 100; email: must be a valid email address

        // Inspect individually
        for _, fe := range check.GetFieldErrors(r) {
            fmt.Printf("Field %q: %s\n", fe.Field, fe.Message)
        }

        // Check specific fields
        if check.HasField(r, "email") {
            fmt.Println("Email validation failed")
        }
    }

    // Fail-fast alternative
    r = check.First(
        check.Required(order.ID, "id"),
        check.UUID(order.ID, "id"),
    )
    // Returns on first error
}
```

## Capabilities

| Category    | Functions                                                                                     |
| ----------- | --------------------------------------------------------------------------------------------- |
| Strings     | `Required`, `MinLen`, `MaxLen`, `Match`, `Prefix`, `Suffix`, `OneOf`, `Alpha`, `Slug`, etc.   |
| Numbers     | `Min`, `Max`, `Between`, `Positive`, `Negative`, `NonZero`, `MultipleOf`, `Percentage`        |
| Comparison  | `Equal`, `NotEqual`, `GreaterThan`, `LessThan`, `EqualField`, `GreaterThanField`              |
| Slices      | `NotEmpty`, `MinItems`, `Unique`, `ContainsAll`, `Each`, `AllSatisfy`, `Subset`               |
| Maps        | `NotEmptyMap`, `HasKey`, `HasKeys`, `OnlyKeys`, `EachKey`, `EachMapValue`, `UniqueValues`     |
| Pointers    | `NotNil`, `Nil`, `NilOr`, `RequiredPtr`, `DefaultOr`, `Deref`                                 |
| Time        | `Before`, `After`, `InPast`, `InFuture`, `BetweenTime`, `WithinDuration`, `NotWeekend`        |
| Formats     | `Email`, `URL`, `UUID`, `IP`, `CIDR`, `Semver`, `E164`, `CreditCard`, `JSON`, `Base64`        |
| Aggregation | `All` (collect all errors), `First` (fail-fast), `Merge` (combine results)                   |

## Validation Tracking

Check tracks which validators were applied to which fields, enabling downstream verification of validation coverage.

```go
r := check.All(
    check.Required(email, "email"),
    check.Email(email, "email"),
    check.Between(age, 13, 120, "age"),
)

// Check if specific validators ran
r.HasValidator("email", "required") // true
r.HasValidator("email", "email")    // true
r.HasValidator("age", "min")        // true (Between reports min and max)
r.HasValidator("age", "max")        // true

// Get all validators for a field
r.ValidatorsFor("email") // []string{"required", "email"}

// Get all validated fields
r.Fields() // []string{"email", "age"}

// Get full tracking map
r.Applied() // map[string][]string{"email": {"required", "email"}, "age": {"min", "max"}}
```

This enables tools to verify that declared validation rules match actual runtime validation.

## Why check?

- **Zero reflection** — validation is function calls, fully visible and debuggable
- **Type-safe generics** — `Min[T]`, `Between[T]`, `Each[T]` catch type errors at compile time
- **Composable** — `All()` collects errors, `First()` fails fast, nest them freely
- **Field-aware errors** — every error knows which field failed and why
- **Validation tracking** — verify which validators ran on which fields
- **No struct tags** — validation rules are code, not strings parsed at runtime
- **Minimal dependencies** — only `golang.org/x/exp` for generic constraints

## Validation as Code

Check enables a pattern: **validation logic as visible, testable code**.

Your validation rules live in methods alongside your types. They're testable, refactorable, and readable. No tag DSL to learn, no reflection overhead, no magic.

```go
// Conditional validation — just Go code
func (o Order) Validate() *check.Result {
    validations := []*check.Validation{
        check.Required(o.ID, "id"),
        check.Positive(o.Total, "total"),
    }

    if o.ShipmentType == "express" {
        validations = append(validations, check.Required(o.ExpressCode, "express_code"))
    }

    return check.All(validations...)
}

// Cross-field validation — just compare values
func (r DateRange) Validate() *check.Result {
    return check.All(
        check.NotZeroTime(r.Start, "start"),
        check.NotZeroTime(r.End, "end"),
        check.GreaterThanField(r.End, r.Start, "end", "start"),
    )
}

// Slice element validation — apply checks to each item
func (c Cart) Validate() *check.Result {
    return check.Merge(
        check.All(check.NotEmpty(c.Items, "items")),
        check.Each(c.Items, func(item Item, i int) *check.Validation {
            field := fmt.Sprintf("items[%d]", i)
            if item.Quantity <= 0 {
                return check.Positive(item.Quantity, field)
            }
            return nil
        }),
    )
}
```

The compiler checks your validation logic. Your IDE can navigate to it. Your tests can exercise it directly.

## Documentation

- [pkg.go.dev reference](https://pkg.go.dev/github.com/zoobzio/check)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License — see [LICENSE](LICENSE) for details.
