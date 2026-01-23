package check

import (
	"errors"
	"regexp"
	"testing"
)

func TestStrBuilder(t *testing.T) {
	t.Run("basic chaining", func(t *testing.T) {
		v := Str("test@example.com", "email").Required().Email().MaxLen(255).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if v.Failed() {
			t.Errorf("expected pass, got error: %v", v.err)
		}
		// Check validators were tracked
		if len(v.validators) != 3 {
			t.Errorf("expected 3 validators, got %d: %v", len(v.validators), v.validators)
		}
	})

	t.Run("collects multiple errors", func(t *testing.T) {
		v := Str("", "email").Required().Email().MinLen(5).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if !v.Failed() {
			t.Error("expected failure")
		}
		// Should have multiple errors collected
		var errs Errors
		if !errors.As(v.err, &errs) {
			t.Fatalf("expected Errors type, got %T", v.err)
		}
		if len(errs) < 2 {
			t.Errorf("expected multiple errors, got %d", len(errs))
		}
	})

	t.Run("works with All", func(t *testing.T) {
		result := All(
			Str("test@example.com", "email").Required().Email().V(),
			Str("password123", "password").Required().MinLen(8).V(),
		)
		if result.Err() != nil {
			t.Errorf("expected pass, got: %v", result.Err())
		}
		// Check applied validators
		if !result.HasValidator("email", "required") {
			t.Error("expected email to have required validator")
		}
		if !result.HasValidator("email", "email") {
			t.Error("expected email to have email validator")
		}
		if !result.HasValidator("password", "min") {
			t.Error("expected password to have min validator")
		}
	})

	t.Run("When conditional", func(t *testing.T) {
		requireStrong := true
		v := Str("weak", "password").
			Required().
			When(requireStrong, func(b *StrBuilder) {
				b.MinLen(12)
			}).
			V()
		if !v.Failed() {
			t.Error("expected failure for weak password with strong requirement")
		}

		// Without condition
		requireStrong = false
		v = Str("weak", "password").
			Required().
			When(requireStrong, func(b *StrBuilder) {
				b.MinLen(12)
			}).
			V()
		if v.Failed() {
			t.Errorf("expected pass without strong requirement, got: %v", v.err)
		}
	})
}

func TestOptStrBuilder(t *testing.T) {
	t.Run("nil pointer skips validation", func(t *testing.T) {
		var name *string
		v := OptStr(name, "name").MaxLen(100).V()
		if v != nil {
			t.Errorf("expected nil for optional nil field, got: %v", v)
		}
	})

	t.Run("non-nil pointer validates", func(t *testing.T) {
		name := "John"
		v := OptStr(&name, "name").MaxLen(100).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("non-nil pointer can fail", func(t *testing.T) {
		name := "This name is way too long for our validation rules"
		v := OptStr(&name, "name").MaxLen(10).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if !v.Failed() {
			t.Error("expected failure for too long name")
		}
	})

	t.Run("works with All", func(t *testing.T) {
		var middleName *string
		lastName := "Doe"

		result := All(
			Str("John", "first_name").Required().V(),
			OptStr(middleName, "middle_name").MaxLen(100).V(),
			OptStr(&lastName, "last_name").MaxLen(100).V(),
		)
		if result.Err() != nil {
			t.Errorf("expected pass, got: %v", result.Err())
		}
	})
}

func TestNumBuilder(t *testing.T) {
	t.Run("basic chaining", func(t *testing.T) {
		v := Num(25, "age").Min(0).Max(120).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("between validation", func(t *testing.T) {
		v := Num(150, "age").Between(0, 120).V()
		if !v.Failed() {
			t.Error("expected failure for age > 120")
		}
	})

	t.Run("works with floats", func(t *testing.T) {
		v := Num(3.14, "pi").Between(3.0, 4.0).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})
}

func TestIntBuilder(t *testing.T) {
	t.Run("integer-specific validators", func(t *testing.T) {
		v := Int(10, "count").Positive().Even().MultipleOf(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("odd validation", func(t *testing.T) {
		v := Int(10, "count").Odd().V()
		if !v.Failed() {
			t.Error("expected failure for even number with Odd()")
		}
	})
}

func TestOptNumBuilder(t *testing.T) {
	t.Run("nil pointer skips validation", func(t *testing.T) {
		var age *int
		v := OptNum(age, "age").Min(0).Max(120).V()
		if v != nil {
			t.Errorf("expected nil for optional nil field, got: %v", v)
		}
	})

	t.Run("non-nil pointer validates", func(t *testing.T) {
		age := 25
		v := OptNum(&age, "age").Min(0).Max(120).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})
}

func TestSliceBuilder(t *testing.T) {
	t.Run("basic slice validation", func(t *testing.T) {
		tags := []string{"go", "rust", "python"}
		v := Slice(tags, "tags").NotEmpty().MaxItems(10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("EachV with auto field names", func(t *testing.T) {
		tags := []string{"go", "rust", "this-tag-is-way-too-long"}
		v := Slice(tags, "tags").
			NotEmpty().
			EachV(func(tag string, field string) *Validation {
				return Str(tag, field).MaxLen(10).V()
			}).
			V()
		if !v.Failed() {
			t.Error("expected failure for long tag")
		}
		// Check the field name is correct
		if v.err == nil {
			t.Fatal("expected error")
		}
		errStr := v.err.Error()
		if errStr == "" {
			t.Error("expected non-empty error")
		}
	})
}

func TestStrSliceBuilder(t *testing.T) {
	t.Run("Each with StrBuilder", func(t *testing.T) {
		tags := []string{"go", "rust", "python"}
		v := StrSlice(tags, "tags").
			NotEmpty().
			MaxItems(10).
			Each(func(b *StrBuilder) {
				b.MaxLen(50).Alpha()
			}).
			V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Each catches element errors", func(t *testing.T) {
		tags := []string{"valid", "also-valid", "has space"}
		v := StrSlice(tags, "tags").
			Each(func(b *StrBuilder) {
				b.NoWhitespace()
			}).
			V()
		if !v.Failed() {
			t.Error("expected failure for tag with space")
		}
	})

	t.Run("AllMaxLen convenience", func(t *testing.T) {
		tags := []string{"go", "rust", "this-is-too-long"}
		v := StrSlice(tags, "tags").AllMaxLen(10).V()
		if !v.Failed() {
			t.Error("expected failure for long tag")
		}
	})

	t.Run("Unique validation", func(t *testing.T) {
		tags := []string{"go", "rust", "go"}
		v := StrSlice(tags, "tags").Unique().V()
		if !v.Failed() {
			t.Error("expected failure for duplicate tags")
		}
	})
}

func TestBuilderIntegration(t *testing.T) {
	// Simulates a real Validate() method
	type CreateUserRequest struct {
		DisplayName *string
		Email       string
		Password    string
		Tags        []string
		Age         int
	}

	validate := func(r *CreateUserRequest) error {
		return All(
			Str(r.Email, "email").Required().Email().MaxLen(255).V(),
			Str(r.Password, "password").Required().MinLen(8).MaxLen(72).V(),
			OptStr(r.DisplayName, "display_name").MaxLen(100).V(),
			Num(r.Age, "age").Min(13).Max(120).V(),
			StrSlice(r.Tags, "tags").MaxItems(10).Each(func(b *StrBuilder) {
				b.NotBlank().MaxLen(50)
			}).V(),
		).Err()
	}

	t.Run("valid request", func(t *testing.T) {
		displayName := "John Doe"
		req := &CreateUserRequest{
			Email:       "john@example.com",
			Password:    "securepassword123",
			DisplayName: &displayName,
			Age:         25,
			Tags:        []string{"developer", "golang"},
		}
		if err := validate(req); err != nil {
			t.Errorf("expected valid, got: %v", err)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		req := &CreateUserRequest{
			Email:    "not-an-email",
			Password: "short",
			Age:      10, // too young
			Tags:     []string{"valid", ""},
		}
		err := validate(req)
		if err == nil {
			t.Error("expected error for invalid request")
		}
	})

	t.Run("optional field nil", func(t *testing.T) {
		req := &CreateUserRequest{
			Email:       "john@example.com",
			Password:    "securepassword123",
			DisplayName: nil, // optional
			Age:         25,
			Tags:        []string{},
		}
		if err := validate(req); err != nil {
			t.Errorf("expected valid with nil optional, got: %v", err)
		}
	})
}

func TestStrBuilderFormatValidators(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *Validation
		wantErr bool
	}{
		{name: "valid email", builder: func() *Validation { return Str("test@example.com", "f").Email().V() }, wantErr: false},
		{name: "invalid email", builder: func() *Validation { return Str("not-email", "f").Email().V() }, wantErr: true},
		{name: "valid url", builder: func() *Validation { return Str("https://example.com", "f").URL().V() }, wantErr: false},
		{name: "invalid url", builder: func() *Validation { return Str("not-a-url", "f").URL().V() }, wantErr: true},
		{name: "valid uuid", builder: func() *Validation { return Str("550e8400-e29b-41d4-a716-446655440000", "f").UUID().V() }, wantErr: false},
		{name: "valid ipv4", builder: func() *Validation { return Str("192.168.1.1", "f").IPv4().V() }, wantErr: false},
		{name: "valid slug", builder: func() *Validation { return Str("my-cool-slug", "f").Slug().V() }, wantErr: false},
		{name: "invalid slug", builder: func() *Validation { return Str("Not A Slug", "f").Slug().V() }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.builder()
			if v == nil {
				t.Fatal("expected validation, got nil")
			}
			if tt.wantErr && !v.Failed() {
				t.Error("expected failure")
			}
			if !tt.wantErr && v.Failed() {
				t.Errorf("expected pass, got: %v", v.err)
			}
		})
	}
}

func TestStrBuilderMatch(t *testing.T) {
	pattern := regexp.MustCompile(`^[A-Z]{3}-\d{4}$`)

	t.Run("matches pattern", func(t *testing.T) {
		v := Str("ABC-1234", "code").Match(pattern).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("does not match pattern", func(t *testing.T) {
		v := Str("abc-1234", "code").Match(pattern).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})
}

func TestCombineEmpty(t *testing.T) {
	// Edge case: builder with no validations
	v := Str("test", "field").V()
	if v != nil {
		t.Errorf("expected nil for empty validations, got: %v", v)
	}
}

func TestStrBuilderStringValidators(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *Validation
		wantErr bool
	}{
		// Len
		{name: "len exact pass", builder: func() *Validation { return Str("hello", "f").Len(5).V() }, wantErr: false},
		{name: "len exact fail", builder: func() *Validation { return Str("hello", "f").Len(3).V() }, wantErr: true},

		// LenBetween
		{name: "len between pass", builder: func() *Validation { return Str("hello", "f").LenBetween(3, 10).V() }, wantErr: false},
		{name: "len between fail short", builder: func() *Validation { return Str("hi", "f").LenBetween(3, 10).V() }, wantErr: true},
		{name: "len between fail long", builder: func() *Validation { return Str("hello world!", "f").LenBetween(3, 10).V() }, wantErr: true},

		// Prefix
		{name: "prefix pass", builder: func() *Validation { return Str("hello world", "f").Prefix("hello").V() }, wantErr: false},
		{name: "prefix fail", builder: func() *Validation { return Str("hello world", "f").Prefix("world").V() }, wantErr: true},

		// Suffix
		{name: "suffix pass", builder: func() *Validation { return Str("hello world", "f").Suffix("world").V() }, wantErr: false},
		{name: "suffix fail", builder: func() *Validation { return Str("hello world", "f").Suffix("hello").V() }, wantErr: true},

		// Contains
		{name: "contains pass", builder: func() *Validation { return Str("hello world", "f").Contains("lo wo").V() }, wantErr: false},
		{name: "contains fail", builder: func() *Validation { return Str("hello world", "f").Contains("xyz").V() }, wantErr: true},

		// NotContains
		{name: "not contains pass", builder: func() *Validation { return Str("hello world", "f").NotContains("xyz").V() }, wantErr: false},
		{name: "not contains fail", builder: func() *Validation { return Str("hello world", "f").NotContains("world").V() }, wantErr: true},

		// OneOf
		{name: "one of pass", builder: func() *Validation { return Str("red", "f").OneOf([]string{"red", "green", "blue"}).V() }, wantErr: false},
		{name: "one of fail", builder: func() *Validation { return Str("yellow", "f").OneOf([]string{"red", "green", "blue"}).V() }, wantErr: true},

		// NotOneOf
		{name: "not one of pass", builder: func() *Validation { return Str("yellow", "f").NotOneOf([]string{"red", "green", "blue"}).V() }, wantErr: false},
		{name: "not one of fail", builder: func() *Validation { return Str("red", "f").NotOneOf([]string{"red", "green", "blue"}).V() }, wantErr: true},

		// AlphaNumeric
		{name: "alphanumeric pass", builder: func() *Validation { return Str("hello123", "f").AlphaNumeric().V() }, wantErr: false},
		{name: "alphanumeric fail", builder: func() *Validation { return Str("hello-123", "f").AlphaNumeric().V() }, wantErr: true},

		// Numeric
		{name: "numeric pass", builder: func() *Validation { return Str("12345", "f").Numeric().V() }, wantErr: false},
		{name: "numeric fail", builder: func() *Validation { return Str("123abc", "f").Numeric().V() }, wantErr: true},

		// AlphaUnicode
		{name: "alpha unicode pass", builder: func() *Validation { return Str("héllo", "f").AlphaUnicode().V() }, wantErr: false},
		{name: "alpha unicode fail", builder: func() *Validation { return Str("hello123", "f").AlphaUnicode().V() }, wantErr: true},

		// AlphaNumericUnicode
		{name: "alphanumeric unicode pass", builder: func() *Validation { return Str("héllo123", "f").AlphaNumericUnicode().V() }, wantErr: false},
		{name: "alphanumeric unicode fail", builder: func() *Validation { return Str("hello-123", "f").AlphaNumericUnicode().V() }, wantErr: true},

		// ASCII
		{name: "ascii pass", builder: func() *Validation { return Str("hello123!@#", "f").ASCII().V() }, wantErr: false},
		{name: "ascii fail", builder: func() *Validation { return Str("héllo", "f").ASCII().V() }, wantErr: true},

		// PrintableASCII
		{name: "printable ascii pass", builder: func() *Validation { return Str("hello 123", "f").PrintableASCII().V() }, wantErr: false},
		{name: "printable ascii fail", builder: func() *Validation { return Str("hello\x00", "f").PrintableASCII().V() }, wantErr: true},

		// LowerCase
		{name: "lowercase pass", builder: func() *Validation { return Str("hello", "f").LowerCase().V() }, wantErr: false},
		{name: "lowercase fail", builder: func() *Validation { return Str("Hello", "f").LowerCase().V() }, wantErr: true},

		// UpperCase
		{name: "uppercase pass", builder: func() *Validation { return Str("HELLO", "f").UpperCase().V() }, wantErr: false},
		{name: "uppercase fail", builder: func() *Validation { return Str("Hello", "f").UpperCase().V() }, wantErr: true},

		// Trimmed
		{name: "trimmed pass", builder: func() *Validation { return Str("hello", "f").Trimmed().V() }, wantErr: false},
		{name: "trimmed fail", builder: func() *Validation { return Str(" hello ", "f").Trimmed().V() }, wantErr: true},

		// SingleLine
		{name: "single line pass", builder: func() *Validation { return Str("hello world", "f").SingleLine().V() }, wantErr: false},
		{name: "single line fail", builder: func() *Validation { return Str("hello\nworld", "f").SingleLine().V() }, wantErr: true},

		// Identifier
		{name: "identifier pass", builder: func() *Validation { return Str("myVar_123", "f").Identifier().V() }, wantErr: false},
		{name: "identifier fail start num", builder: func() *Validation { return Str("123var", "f").Identifier().V() }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.builder()
			if v == nil {
				t.Fatal("expected validation, got nil")
			}
			if tt.wantErr && !v.Failed() {
				t.Error("expected failure")
			}
			if !tt.wantErr && v.Failed() {
				t.Errorf("expected pass, got: %v", v.err)
			}
		})
	}
}

func TestStrBuilderNotMatch(t *testing.T) {
	pattern := regexp.MustCompile(`\d+`)

	t.Run("not match pass", func(t *testing.T) {
		v := Str("hello", "f").NotMatch(pattern).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("not match fail", func(t *testing.T) {
		v := Str("hello123", "f").NotMatch(pattern).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})
}

func TestStrBuilderFormatValidatorsExtended(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *Validation
		wantErr bool
	}{
		// URLWithScheme
		{name: "url with scheme pass", builder: func() *Validation { return Str("https://example.com", "f").URLWithScheme([]string{"https"}).V() }, wantErr: false},
		{name: "url with scheme fail", builder: func() *Validation { return Str("http://example.com", "f").URLWithScheme([]string{"https"}).V() }, wantErr: true},

		// HTTPOrHTTPS
		{name: "http or https pass http", builder: func() *Validation { return Str("http://example.com", "f").HTTPOrHTTPS().V() }, wantErr: false},
		{name: "http or https pass https", builder: func() *Validation { return Str("https://example.com", "f").HTTPOrHTTPS().V() }, wantErr: false},
		{name: "http or https fail", builder: func() *Validation { return Str("ftp://example.com", "f").HTTPOrHTTPS().V() }, wantErr: true},

		// UUID4
		{name: "uuid4 pass", builder: func() *Validation { return Str("a3bb189e-8bf9-4e15-9c87-0b8d0c5b36f3", "f").UUID4().V() }, wantErr: false},

		// IP
		{name: "ip pass v4", builder: func() *Validation { return Str("192.168.1.1", "f").IP().V() }, wantErr: false},
		{name: "ip pass v6", builder: func() *Validation { return Str("::1", "f").IP().V() }, wantErr: false},
		{name: "ip fail", builder: func() *Validation { return Str("not-an-ip", "f").IP().V() }, wantErr: true},

		// IPv6
		{name: "ipv6 pass", builder: func() *Validation { return Str("2001:0db8:85a3:0000:0000:8a2e:0370:7334", "f").IPv6().V() }, wantErr: false},
		{name: "ipv6 fail", builder: func() *Validation { return Str("192.168.1.1", "f").IPv6().V() }, wantErr: true},

		// CIDR
		{name: "cidr pass", builder: func() *Validation { return Str("192.168.1.0/24", "f").CIDR().V() }, wantErr: false},
		{name: "cidr fail", builder: func() *Validation { return Str("192.168.1.1", "f").CIDR().V() }, wantErr: true},

		// MAC
		{name: "mac pass", builder: func() *Validation { return Str("00:1A:2B:3C:4D:5E", "f").MAC().V() }, wantErr: false},
		{name: "mac fail", builder: func() *Validation { return Str("not-a-mac", "f").MAC().V() }, wantErr: true},

		// Hostname
		{name: "hostname pass", builder: func() *Validation { return Str("example.com", "f").Hostname().V() }, wantErr: false},
		{name: "hostname fail", builder: func() *Validation { return Str("not a hostname!", "f").Hostname().V() }, wantErr: true},

		// Port
		{name: "port pass", builder: func() *Validation { return Str("8080", "f").Port().V() }, wantErr: false},
		{name: "port fail", builder: func() *Validation { return Str("99999", "f").Port().V() }, wantErr: true},

		// HostPort
		{name: "host port pass", builder: func() *Validation { return Str("example.com:8080", "f").HostPort().V() }, wantErr: false},
		{name: "host port fail", builder: func() *Validation { return Str("example.com", "f").HostPort().V() }, wantErr: true},

		// HexColor
		{name: "hex color pass short", builder: func() *Validation { return Str("#fff", "f").HexColor().V() }, wantErr: false},
		{name: "hex color pass full", builder: func() *Validation { return Str("#ffffff", "f").HexColor().V() }, wantErr: false},
		{name: "hex color fail", builder: func() *Validation { return Str("red", "f").HexColor().V() }, wantErr: true},

		// HexColorFull
		{name: "hex color full pass", builder: func() *Validation { return Str("#ffffff", "f").HexColorFull().V() }, wantErr: false},
		{name: "hex color full fail short", builder: func() *Validation { return Str("#fff", "f").HexColorFull().V() }, wantErr: true},

		// Base64
		{name: "base64 pass", builder: func() *Validation { return Str("SGVsbG8gV29ybGQ=", "f").Base64().V() }, wantErr: false},
		{name: "base64 fail", builder: func() *Validation { return Str("not base64!", "f").Base64().V() }, wantErr: true},

		// Base64URL
		{name: "base64url pass", builder: func() *Validation { return Str("dGVzdA==", "f").Base64URL().V() }, wantErr: false},

		// JSON
		{name: "json pass", builder: func() *Validation { return Str(`{"key":"value"}`, "f").JSON().V() }, wantErr: false},
		{name: "json fail", builder: func() *Validation { return Str("not json", "f").JSON().V() }, wantErr: true},

		// Semver
		{name: "semver pass", builder: func() *Validation { return Str("1.2.3", "f").Semver().V() }, wantErr: false},
		{name: "semver fail", builder: func() *Validation { return Str("1.2", "f").Semver().V() }, wantErr: true},

		// E164
		{name: "e164 pass", builder: func() *Validation { return Str("+14155552671", "f").E164().V() }, wantErr: false},
		{name: "e164 fail", builder: func() *Validation { return Str("555-1234", "f").E164().V() }, wantErr: true},

		// CreditCard
		{name: "credit card pass", builder: func() *Validation { return Str("4111111111111111", "f").CreditCard().V() }, wantErr: false},
		{name: "credit card fail", builder: func() *Validation { return Str("1234567890", "f").CreditCard().V() }, wantErr: true},

		// Latitude
		{name: "latitude pass", builder: func() *Validation { return Str("37.7749", "f").Latitude().V() }, wantErr: false},
		{name: "latitude fail", builder: func() *Validation { return Str("91.0", "f").Latitude().V() }, wantErr: true},

		// Longitude
		{name: "longitude pass", builder: func() *Validation { return Str("-122.4194", "f").Longitude().V() }, wantErr: false},
		{name: "longitude fail", builder: func() *Validation { return Str("181.0", "f").Longitude().V() }, wantErr: true},

		// CountryCode2
		{name: "country code 2 pass", builder: func() *Validation { return Str("US", "f").CountryCode2().V() }, wantErr: false},
		{name: "country code 2 fail", builder: func() *Validation { return Str("USA", "f").CountryCode2().V() }, wantErr: true},

		// CountryCode3
		{name: "country code 3 pass", builder: func() *Validation { return Str("USA", "f").CountryCode3().V() }, wantErr: false},
		{name: "country code 3 fail", builder: func() *Validation { return Str("US", "f").CountryCode3().V() }, wantErr: true},

		// LanguageCode
		{name: "language code pass", builder: func() *Validation { return Str("en", "f").LanguageCode().V() }, wantErr: false},
		{name: "language code fail", builder: func() *Validation { return Str("english", "f").LanguageCode().V() }, wantErr: true},

		// CurrencyCode
		{name: "currency code pass", builder: func() *Validation { return Str("USD", "f").CurrencyCode().V() }, wantErr: false},
		{name: "currency code fail", builder: func() *Validation { return Str("DOLLAR", "f").CurrencyCode().V() }, wantErr: true},

		// Hex
		{name: "hex pass", builder: func() *Validation { return Str("deadbeef", "f").Hex().V() }, wantErr: false},
		{name: "hex fail", builder: func() *Validation { return Str("not hex!", "f").Hex().V() }, wantErr: true},

		// DataURI
		{name: "data uri pass", builder: func() *Validation { return Str("data:text/plain;base64,SGVsbG8=", "f").DataURI().V() }, wantErr: false},
		{name: "data uri fail", builder: func() *Validation { return Str("not a data uri", "f").DataURI().V() }, wantErr: true},

		// FilePath
		{name: "file path pass", builder: func() *Validation { return Str("/path/to/file.txt", "f").FilePath().V() }, wantErr: false},

		// UnixPath
		{name: "unix path pass", builder: func() *Validation { return Str("/usr/local/bin", "f").UnixPath().V() }, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.builder()
			if v == nil {
				t.Fatal("expected validation, got nil")
			}
			if tt.wantErr && !v.Failed() {
				t.Error("expected failure")
			}
			if !tt.wantErr && v.Failed() {
				t.Errorf("expected pass, got: %v", v.err)
			}
		})
	}
}

func TestOptStrBuilderMethods(t *testing.T) {
	t.Run("When conditional with value", func(t *testing.T) {
		val := "test"
		v := OptStr(&val, "f").When(true, func(b *OptStrBuilder) {
			b.MinLen(10)
		}).V()
		if !v.Failed() {
			t.Error("expected failure for short string with condition")
		}
	})

	t.Run("When conditional without value", func(t *testing.T) {
		var val *string
		v := OptStr(val, "f").When(true, func(b *OptStrBuilder) {
			b.MinLen(10)
		}).V()
		if v != nil {
			t.Errorf("expected nil for nil value, got: %v", v)
		}
	})

	t.Run("MinLen", func(t *testing.T) {
		val := "hi"
		v := OptStr(&val, "f").MinLen(5).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("Len", func(t *testing.T) {
		val := "hello"
		v := OptStr(&val, "f").Len(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("LenBetween", func(t *testing.T) {
		val := "hello"
		v := OptStr(&val, "f").LenBetween(3, 10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Match", func(t *testing.T) {
		pattern := regexp.MustCompile(`^\d+$`)
		val := "123"
		v := OptStr(&val, "f").Match(pattern).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("NotMatch", func(t *testing.T) {
		pattern := regexp.MustCompile(`\d`)
		val := "hello"
		v := OptStr(&val, "f").NotMatch(pattern).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("string content validators", func(t *testing.T) {
		val := "hello"
		tests := []struct {
			name    string
			builder func() *Validation
			wantErr bool
		}{
			{"Prefix", func() *Validation { return OptStr(&val, "f").Prefix("hel").V() }, false},
			{"Suffix", func() *Validation { return OptStr(&val, "f").Suffix("llo").V() }, false},
			{"Contains", func() *Validation { return OptStr(&val, "f").Contains("ell").V() }, false},
			{"NotContains", func() *Validation { return OptStr(&val, "f").NotContains("xyz").V() }, false},
			{"OneOf", func() *Validation { return OptStr(&val, "f").OneOf([]string{"hello", "world"}).V() }, false},
			{"NotOneOf", func() *Validation { return OptStr(&val, "f").NotOneOf([]string{"foo", "bar"}).V() }, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				v := tt.builder()
				if v == nil {
					t.Fatal("expected validation, got nil")
				}
				if tt.wantErr && !v.Failed() {
					t.Error("expected failure")
				}
				if !tt.wantErr && v.Failed() {
					t.Errorf("expected pass, got: %v", v.err)
				}
			})
		}
	})

	t.Run("character validators", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			builder func(s *string) *Validation
			wantErr bool
		}{
			{"Alpha", "hello", func(s *string) *Validation { return OptStr(s, "f").Alpha().V() }, false},
			{"AlphaNumeric", "hello123", func(s *string) *Validation { return OptStr(s, "f").AlphaNumeric().V() }, false},
			{"Numeric", "123", func(s *string) *Validation { return OptStr(s, "f").Numeric().V() }, false},
			{"LowerCase", "hello", func(s *string) *Validation { return OptStr(s, "f").LowerCase().V() }, false},
			{"UpperCase", "HELLO", func(s *string) *Validation { return OptStr(s, "f").UpperCase().V() }, false},
			{"Trimmed", "hello", func(s *string) *Validation { return OptStr(s, "f").Trimmed().V() }, false},
			{"SingleLine", "hello", func(s *string) *Validation { return OptStr(s, "f").SingleLine().V() }, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				v := tt.builder(&tt.value)
				if v == nil {
					t.Fatal("expected validation, got nil")
				}
				if tt.wantErr && !v.Failed() {
					t.Error("expected failure")
				}
				if !tt.wantErr && v.Failed() {
					t.Errorf("expected pass, got: %v", v.err)
				}
			})
		}
	})

	t.Run("format validators", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			builder func(s *string) *Validation
			wantErr bool
		}{
			{"Slug", "my-slug", func(s *string) *Validation { return OptStr(s, "f").Slug().V() }, false},
			{"Email", "test@example.com", func(s *string) *Validation { return OptStr(s, "f").Email().V() }, false},
			{"URL", "https://example.com", func(s *string) *Validation { return OptStr(s, "f").URL().V() }, false},
			{"UUID", "550e8400-e29b-41d4-a716-446655440000", func(s *string) *Validation { return OptStr(s, "f").UUID().V() }, false},
			{"UUID4", "a3bb189e-8bf9-4e15-9c87-0b8d0c5b36f3", func(s *string) *Validation { return OptStr(s, "f").UUID4().V() }, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				v := tt.builder(&tt.value)
				if v == nil {
					t.Fatal("expected validation, got nil")
				}
				if tt.wantErr && !v.Failed() {
					t.Error("expected failure")
				}
				if !tt.wantErr && v.Failed() {
					t.Errorf("expected pass, got: %v", v.err)
				}
			})
		}
	})
}

func TestNumBuilderExtended(t *testing.T) {
	t.Run("When conditional", func(t *testing.T) {
		v := Num(5, "f").When(true, func(b *NumBuilder[int]) {
			b.Min(10)
		}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("BetweenExclusive", func(t *testing.T) {
		v := Num(5, "f").BetweenExclusive(0, 10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Num(0, "f").BetweenExclusive(0, 10).V()
		if !v.Failed() {
			t.Error("expected failure at boundary")
		}
	})

	t.Run("GreaterThan", func(t *testing.T) {
		v := Num(10, "f").GreaterThan(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Num(5, "f").GreaterThan(5).V()
		if !v.Failed() {
			t.Error("expected failure at equal")
		}
	})

	t.Run("LessThan", func(t *testing.T) {
		v := Num(5, "f").LessThan(10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Num(10, "f").LessThan(10).V()
		if !v.Failed() {
			t.Error("expected failure at equal")
		}
	})

	t.Run("GreaterThanOrEqual", func(t *testing.T) {
		v := Num(5, "f").GreaterThanOrEqual(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Num(4, "f").GreaterThanOrEqual(5).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("LessThanOrEqual", func(t *testing.T) {
		v := Num(5, "f").LessThanOrEqual(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Num(6, "f").LessThanOrEqual(5).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("OneOfValues", func(t *testing.T) {
		v := Num(5, "f").OneOfValues([]int{1, 3, 5, 7}).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Num(2, "f").OneOfValues([]int{1, 3, 5, 7}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("NotOneOfValues", func(t *testing.T) {
		v := Num(2, "f").NotOneOfValues([]int{1, 3, 5, 7}).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Num(5, "f").NotOneOfValues([]int{1, 3, 5, 7}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})
}

func TestIntBuilderExtended(t *testing.T) {
	t.Run("When conditional", func(t *testing.T) {
		v := Int(5, "f").When(true, func(b *IntBuilder[int]) {
			b.Min(10)
		}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("Min", func(t *testing.T) {
		v := Int(10, "f").Min(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Max", func(t *testing.T) {
		v := Int(5, "f").Max(10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Between", func(t *testing.T) {
		v := Int(5, "f").Between(0, 10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Negative", func(t *testing.T) {
		v := Int(-5, "f").Negative().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Int(5, "f").Negative().V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("NonNegative", func(t *testing.T) {
		v := Int(0, "f").NonNegative().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Int(-1, "f").NonNegative().V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("NonPositive", func(t *testing.T) {
		v := Int(0, "f").NonPositive().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Int(1, "f").NonPositive().V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("Zero", func(t *testing.T) {
		v := Int(0, "f").Zero().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Int(1, "f").Zero().V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("NonZero", func(t *testing.T) {
		v := Int(5, "f").NonZero().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Int(0, "f").NonZero().V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})
}

func TestOptNumBuilderExtended(t *testing.T) {
	t.Run("When conditional", func(t *testing.T) {
		val := 5
		v := OptNum(&val, "f").When(true, func(b *OptNumBuilder[int]) {
			b.Min(10)
		}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("Between", func(t *testing.T) {
		val := 5
		v := OptNum(&val, "f").Between(0, 10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("GreaterThan", func(t *testing.T) {
		val := 10
		v := OptNum(&val, "f").GreaterThan(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("LessThan", func(t *testing.T) {
		val := 5
		v := OptNum(&val, "f").LessThan(10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})
}

func TestSliceBuilderExtended(t *testing.T) {
	t.Run("When conditional", func(t *testing.T) {
		items := []int{1, 2, 3}
		v := Slice(items, "f").When(true, func(b *SliceBuilder[int]) {
			b.MinItems(5)
		}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		var items []int
		v := Slice(items, "f").Empty().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		items = []int{1}
		v = Slice(items, "f").Empty().V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("MinItems", func(t *testing.T) {
		items := []int{1, 2, 3}
		v := Slice(items, "f").MinItems(2).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Slice(items, "f").MinItems(5).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("ExactItems", func(t *testing.T) {
		items := []int{1, 2, 3}
		v := Slice(items, "f").ExactItems(3).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		v = Slice(items, "f").ExactItems(5).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("ItemsBetween", func(t *testing.T) {
		items := []int{1, 2, 3}
		v := Slice(items, "f").ItemsBetween(1, 5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Each", func(t *testing.T) {
		items := []int{1, 2, 3}
		count := 0
		_ = Slice(items, "f").Each(func(_ int, _ string) *Validation {
			count++
			return nil
		}).V()
		if count != 3 {
			t.Errorf("expected 3 iterations, got %d", count)
		}
	})

	t.Run("Each collects validations", func(t *testing.T) {
		items := []int{1, 2, 100}
		v := Slice(items, "f").Each(func(item int, field string) *Validation {
			return Num(item, field).Max(50).V()
		}).V()
		if !v.Failed() {
			t.Error("expected failure for item > 50")
		}
	})
}

func TestStrSliceBuilderExtended(t *testing.T) {
	t.Run("When conditional", func(t *testing.T) {
		items := []string{"a", "b"}
		v := StrSlice(items, "f").When(true, func(b *StrSliceBuilder) {
			b.MinItems(5)
		}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("MinItems", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		v := StrSlice(items, "f").MinItems(2).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("ItemsBetween", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		v := StrSlice(items, "f").ItemsBetween(1, 5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("AllMinLen", func(t *testing.T) {
		items := []string{"hello", "world"}
		v := StrSlice(items, "f").AllMinLen(3).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		items = []string{"hi", "world"}
		v = StrSlice(items, "f").AllMinLen(3).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("AllNotBlank", func(t *testing.T) {
		items := []string{"hello", "world"}
		v := StrSlice(items, "f").AllNotBlank().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		items = []string{"hello", "   "}
		v = StrSlice(items, "f").AllNotBlank().V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})
}

func TestOptIntBuilder(t *testing.T) {
	t.Run("nil skips validation", func(t *testing.T) {
		var val *int
		v := OptInt(val, "f").Min(0).V()
		if v != nil {
			t.Errorf("expected nil, got: %v", v)
		}
	})

	t.Run("non-nil validates", func(t *testing.T) {
		val := 10
		v := OptInt(&val, "f").Min(0).Max(100).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("When conditional", func(t *testing.T) {
		val := 5
		v := OptInt(&val, "f").When(true, func(b *OptIntBuilder[int]) {
			b.Min(10)
		}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("Between", func(t *testing.T) {
		val := 5
		v := OptInt(&val, "f").Between(0, 10).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Positive", func(t *testing.T) {
		val := 5
		v := OptInt(&val, "f").Positive().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}

		val = -1
		v = OptInt(&val, "f").Positive().V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("NonNegative", func(t *testing.T) {
		val := 0
		v := OptInt(&val, "f").NonNegative().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("NonZero", func(t *testing.T) {
		val := 5
		v := OptInt(&val, "f").NonZero().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("MultipleOf", func(t *testing.T) {
		val := 10
		v := OptInt(&val, "f").MultipleOf(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Even", func(t *testing.T) {
		val := 10
		v := OptInt(&val, "f").Even().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Odd", func(t *testing.T) {
		val := 11
		v := OptInt(&val, "f").Odd().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})
}

func TestOptSliceBuilder(t *testing.T) {
	t.Run("nil skips validation", func(t *testing.T) {
		var items *[]int
		v := OptSlice(items, "f").MinItems(1).V()
		if v != nil {
			t.Errorf("expected nil, got: %v", v)
		}
	})

	t.Run("non-nil validates", func(t *testing.T) {
		items := []int{1, 2, 3}
		v := OptSlice(&items, "f").NotEmpty().MaxItems(10).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("When conditional", func(t *testing.T) {
		items := []int{1, 2}
		v := OptSlice(&items, "f").When(true, func(b *OptSliceBuilder[int]) {
			b.MinItems(5)
		}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("NotEmpty", func(t *testing.T) {
		items := []int{1}
		v := OptSlice(&items, "f").NotEmpty().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("MinItems", func(t *testing.T) {
		items := []int{1, 2, 3}
		v := OptSlice(&items, "f").MinItems(2).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("MaxItems", func(t *testing.T) {
		items := []int{1, 2, 3}
		v := OptSlice(&items, "f").MaxItems(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("ItemsBetween", func(t *testing.T) {
		items := []int{1, 2, 3}
		v := OptSlice(&items, "f").ItemsBetween(1, 5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("EachV", func(t *testing.T) {
		items := []int{1, 2, 100}
		v := OptSlice(&items, "f").EachV(func(item int, field string) *Validation {
			return Num(item, field).Max(50).V()
		}).V()
		if !v.Failed() {
			t.Error("expected failure for item > 50")
		}
	})
}

func TestOptStrSliceBuilder(t *testing.T) {
	t.Run("nil skips validation", func(t *testing.T) {
		var items *[]string
		v := OptStrSlice(items, "f").MinItems(1).V()
		if v != nil {
			t.Errorf("expected nil, got: %v", v)
		}
	})

	t.Run("non-nil validates", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		v := OptStrSlice(&items, "f").NotEmpty().MaxItems(10).V()
		if v == nil {
			t.Fatal("expected validation, got nil")
		}
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("When conditional", func(t *testing.T) {
		items := []string{"a", "b"}
		v := OptStrSlice(&items, "f").When(true, func(b *OptStrSliceBuilder) {
			b.MinItems(5)
		}).V()
		if !v.Failed() {
			t.Error("expected failure")
		}
	})

	t.Run("NotEmpty", func(t *testing.T) {
		items := []string{"a"}
		v := OptStrSlice(&items, "f").NotEmpty().V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("MinItems", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		v := OptStrSlice(&items, "f").MinItems(2).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("MaxItems", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		v := OptStrSlice(&items, "f").MaxItems(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		items := []string{"a", "b", "a"}
		v := OptStrSlice(&items, "f").Unique().V()
		if !v.Failed() {
			t.Error("expected failure for duplicates")
		}
	})

	t.Run("Each", func(t *testing.T) {
		items := []string{"hello", "world"}
		v := OptStrSlice(&items, "f").Each(func(b *StrBuilder) {
			b.MinLen(3)
		}).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("AllMaxLen", func(t *testing.T) {
		items := []string{"hi", "yo"}
		v := OptStrSlice(&items, "f").AllMaxLen(5).V()
		if v.Failed() {
			t.Errorf("expected pass, got: %v", v.err)
		}
	})

	t.Run("AllNotBlank", func(t *testing.T) {
		items := []string{"hello", "   "}
		v := OptStrSlice(&items, "f").AllNotBlank().V()
		if !v.Failed() {
			t.Error("expected failure for blank item")
		}
	})
}
