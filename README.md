# check

[![CI Status](https://github.com/zoobzio/check/workflows/CI/badge.svg)](https://github.com/zoobzio/check/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/zoobzio/check/graph/badge.svg?branch=main)](https://codecov.io/gh/zoobzio/check)
[![Go Report Card](https://goreportcard.com/badge/github.com/zoobzio/check)](https://goreportcard.com/report/github.com/zoobzio/check)
[![CodeQL](https://github.com/zoobzio/check/workflows/CodeQL/badge.svg)](https://github.com/zoobzio/check/security/code-scanning)
[![Go Reference](https://pkg.go.dev/badge/github.com/zoobzio/check.svg)](https://pkg.go.dev/github.com/zoobzio/check)
[![License](https://img.shields.io/github/license/zoobzio/check)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod-go-version/zoobzio/check)](go.mod)
[![Release](https://img.shields.io/github/v/release/zoobzio/check)](https://github.com/zoobzio/check/releases)

Fluent validation for Go with struct tag verification.

## Usage

```go
type User struct {
    Email    string  `json:"email" validate:"required,email"`
    Password string  `json:"password" validate:"required,min=8"`
    Name     *string `json:"name" validate:"omitempty,max=100"`
    Age      int     `json:"age" validate:"min=13,max=120"`
}

func (u *User) Validate() error {
    return check.Check[User](
        check.Str(u.Email, "email").Required().Email().MaxLen(255).V(),
        check.Str(u.Password, "password").Required().MinLen(8).V(),
        check.OptStr(u.Name, "name").MaxLen(100).V(),
        check.Num(u.Age, "age").Between(13, 120).V(),
    ).Err()
}
```

`Check[T]` validates your fields and verifies that every field with a `validate` tag was actually checked. Forget a field? You'll know:

```
password: tagged but not validated (validate: required,min=8)
```

No magic, no reflection at validation time—just functions that return validation results.

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

## Fluent Builders

Chain validators naturally:

```go
// Strings
check.Str(email, "email").Required().Email().MaxLen(255).V()

// Optional strings (nil skips validation)
check.OptStr(name, "name").MaxLen(100).V()

// Numbers
check.Num(age, "age").Between(13, 120).V()

// Integers (adds Even, Odd, MultipleOf)
check.Int(count, "count").Positive().Even().V()

// Slices with auto-generated field names
check.StrSlice(tags, "tags").NotEmpty().MaxItems(10).Each(func(b *check.StrBuilder) {
    b.MaxLen(50)  // Validates tags[0], tags[1], etc.
}).V()
```

Conditional validation with `.When()`:

```go
check.Str(password, "password").
    Required().
    When(requireStrong, func(b *check.StrBuilder) {
        b.MinLen(12).Match(complexityRegex)
    }).V()
```

## Direct Functions

Use validators directly when you don't need the fluent API:

```go
check.All(
    check.Required(email, "email"),
    check.Email(email, "email"),
    check.Between(age, 13, 120, "age"),
)
```

## Capabilities

| Category    | Builders                                           | Functions                                                                                   |
| ----------- | -------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| Strings     | `Str`, `OptStr`                                    | `Required`, `MinLen`, `MaxLen`, `Match`, `Prefix`, `Suffix`, `OneOf`, `Alpha`, `Slug`, etc. |
| Numbers     | `Num`, `OptNum`, `Int`, `OptInt`                   | `Min`, `Max`, `Between`, `Positive`, `Negative`, `NonZero`, `MultipleOf`, `Percentage`      |
| Slices      | `Slice`, `OptSlice`, `StrSlice`, `OptStrSlice`     | `NotEmpty`, `MinItems`, `Unique`, `ContainsAll`, `Each`, `AllSatisfy`, `Subset`             |
| Formats     | (via `Str` methods)                                | `Email`, `URL`, `UUID`, `IP`, `CIDR`, `Semver`, `E164`, `CreditCard`, `JSON`, `Base64`      |
| Comparison  | —                                                  | `Equal`, `NotEqual`, `GreaterThan`, `LessThan`, `EqualField`, `GreaterThanField`            |
| Maps        | —                                                  | `NotEmptyMap`, `HasKey`, `HasKeys`, `OnlyKeys`, `EachKey`, `EachMapValue`, `UniqueValues`   |
| Pointers    | —                                                  | `NotNil`, `Nil`, `NilOr`, `RequiredPtr`, `DefaultOr`, `Deref`                               |
| Time        | —                                                  | `Before`, `After`, `InPast`, `InFuture`, `BetweenTime`, `WithinDuration`, `NotWeekend`      |
| Aggregation | —                                                  | `All` (collect all errors), `First` (fail-fast), `Merge`, `Check[T]` (with tag verification)|

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

- **Fluent API** — chain validators, reduce boilerplate
- **Tag verification** — `Check[T]` catches forgotten fields at runtime
- **Zero reflection** — validation is function calls, fully visible and debuggable
- **Type-safe generics** — `Min[T]`, `Between[T]`, `Each[T]` catch type errors at compile time
- **Composable** — `All()` collects errors, `First()` fails fast, nest them freely
- **Field-aware errors** — every error knows which field failed and why
- **Validation tracking** — verify which validators ran on which fields

## Validation as Code

Check enables a pattern: **validation logic as visible, testable code**.

Your validation rules live in methods alongside your types. They're testable, refactorable, and readable. No tag DSL to learn, no reflection overhead, no magic.

```go
// Conditional validation — just Go code
func (o Order) Validate() error {
    validations := []*check.Validation{
        check.Str(o.ID, "id").Required().UUID().V(),
        check.Num(o.Total, "total").Positive().V(),
    }

    if o.ShipmentType == "express" {
        validations = append(validations,
            check.Str(o.ExpressCode, "express_code").Required().V(),
        )
    }

    return check.All(validations...).Err()
}

// Cross-field validation — just compare values
func (r DateRange) Validate() error {
    return check.All(
        check.NotZeroTime(r.Start, "start"),
        check.NotZeroTime(r.End, "end"),
        check.GreaterThanField(r.End, r.Start, "end", "start"),
    ).Err()
}

// Slice element validation — auto-generated field names
func (c Cart) Validate() error {
    return check.All(
        check.StrSlice(c.ItemIDs, "items").NotEmpty().Each(func(b *check.StrBuilder) {
            b.Required().UUID()
        }).V(),
    ).Err()
}
```

The compiler checks your validation logic. Your IDE can navigate to it. Your tests can exercise it directly.

## Documentation

- [pkg.go.dev reference](https://pkg.go.dev/github.com/zoobzio/check)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License — see [LICENSE](LICENSE) for details.
