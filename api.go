// Package check provides zero-reflection validation primitives for Go.
//
// # Overview
//
// check offers explicit validation functions with no struct tags or reflection.
// Validators are composable via [All] for collecting multiple errors, or [First]
// for fail-fast behavior.
//
// # Basic Usage
//
//	func (r *Request) Validate() error {
//	    return check.All(
//	        check.Required(r.Email, "email"),
//	        check.Email(r.Email, "email"),
//	        check.MinLen(r.Password, 8, "password"),
//	    )
//	}
//
// # Error Handling
//
// All validators return [*FieldError] on failure, which includes the field name
// and a descriptive message. The [All] function collects errors into an [Errors]
// slice, while [First] returns only the first error encountered.
//
//	err := check.All(
//	    check.Required(name, "name"),
//	    check.Email(email, "email"),
//	)
//	if err != nil {
//	    for _, fe := range check.GetFieldErrors(err) {
//	        fmt.Printf("%s: %s\n", fe.Field, fe.Message)
//	    }
//	}
//
// # Validator Categories
//
// String validators: [Required], [MinLen], [MaxLen], [Len], [Match], [OneOf],
// [Alpha], [AlphaNumeric], [Slug], [Identifier], and more.
//
// Format validators: [Email], [URL], [UUID], [IP], [IPv4], [IPv6], [CIDR],
// [MAC], [HexColor], [JSON], [Semver], [E164], [CreditCard], and more.
//
// Numeric validators: [Min], [Max], [Between], [Positive], [Negative],
// [NonZero], [MultipleOf], [Even], [Odd], and more.
//
// Slice validators: [NotEmpty], [MinItems], [MaxItems], [Unique], [Each],
// [ContainsAll], [ContainsAny], [Subset], and more.
//
// Map validators: [NotEmptyMap], [HasKey], [HasKeys], [OnlyKeys],
// [UniqueValues], and more.
//
// Pointer validators: [NotNil], [NilOr], [RequiredPtr], [Deref], and more.
//
// Time validators: [Before], [After], [BetweenTime], [WithinDuration],
// [SameDay], [Weekday], [NotWeekend], and more.
//
// Comparison validators: [Equal], [NotEqual], [EqualField], [GreaterThanField],
// and more.
//
// # Optional Field Validation
//
// Use [NilOr] to validate pointer fields only when present:
//
//	check.NilOr(r.MiddleName, func(v string) error {
//	    return check.MaxLen(v, 100, "middle_name")
//	})
//
// # Slice Element Validation
//
// Use [Each] to validate every element in a slice:
//
//	check.Each(r.Tags, func(tag string, i int) error {
//	    return check.All(
//	        check.Required(tag, fmt.Sprintf("tags[%d]", i)),
//	        check.MaxLen(tag, 50, fmt.Sprintf("tags[%d]", i)),
//	    )
//	}, "tags")
package check
