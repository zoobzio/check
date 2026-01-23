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
