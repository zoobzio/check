package check

import (
	"fmt"
	"testing"
)

func TestNotEmpty(t *testing.T) {
	tests := []struct {
		value   []int
		wantErr bool
	}{
		{[]int{1, 2, 3}, false},
		{[]int{1}, false},
		{[]int{}, true},
		{nil, true},
	}
	for _, tt := range tests {
		v := NotEmpty(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NotEmpty(%v) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		value   []int
		wantErr bool
	}{
		{[]int{}, false},
		{nil, false},
		{[]int{1}, true},
	}
	for _, tt := range tests {
		v := Empty(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Empty(%v) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestMinItems(t *testing.T) {
	tests := []struct {
		value   []int
		min     int
		wantErr bool
	}{
		{[]int{1, 2, 3}, 2, false},
		{[]int{1, 2}, 2, false},
		{[]int{1}, 2, true},
	}
	for _, tt := range tests {
		v := MinItems(tt.value, tt.min, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("MinItems(%v, %d) failed = %v, wantErr %v", tt.value, tt.min, v.Failed(), tt.wantErr)
		}
	}
}

func TestMaxItems(t *testing.T) {
	tests := []struct {
		value   []int
		max     int
		wantErr bool
	}{
		{[]int{1}, 2, false},
		{[]int{1, 2}, 2, false},
		{[]int{1, 2, 3}, 2, true},
	}
	for _, tt := range tests {
		v := MaxItems(tt.value, tt.max, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("MaxItems(%v, %d) failed = %v, wantErr %v", tt.value, tt.max, v.Failed(), tt.wantErr)
		}
	}
}

func TestExactItems(t *testing.T) {
	tests := []struct {
		value   []int
		count   int
		wantErr bool
	}{
		{[]int{1, 2}, 2, false},
		{[]int{1}, 2, true},
		{[]int{1, 2, 3}, 2, true},
	}
	for _, tt := range tests {
		v := ExactItems(tt.value, tt.count, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("ExactItems(%v, %d) failed = %v, wantErr %v", tt.value, tt.count, v.Failed(), tt.wantErr)
		}
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		value   []int
		wantErr bool
	}{
		{[]int{1, 2, 3}, false},
		{[]int{1}, false},
		{[]int{}, false},
		{[]int{1, 2, 1}, true},
	}
	for _, tt := range tests {
		v := Unique(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Unique(%v) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestSliceContains(t *testing.T) {
	tests := []struct {
		value   []int
		elem    int
		wantErr bool
	}{
		{[]int{1, 2, 3}, 2, false},
		{[]int{1, 2, 3}, 4, true},
		{[]int{}, 1, true},
	}
	for _, tt := range tests {
		v := SliceContains(tt.value, tt.elem, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("SliceContains(%v, %d) failed = %v, wantErr %v", tt.value, tt.elem, v.Failed(), tt.wantErr)
		}
	}
}

func TestSliceNotContains(t *testing.T) {
	tests := []struct {
		value   []int
		elem    int
		wantErr bool
	}{
		{[]int{1, 2, 3}, 4, false},
		{[]int{}, 1, false},
		{[]int{1, 2, 3}, 2, true},
	}
	for _, tt := range tests {
		v := SliceNotContains(tt.value, tt.elem, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("SliceNotContains(%v, %d) failed = %v, wantErr %v", tt.value, tt.elem, v.Failed(), tt.wantErr)
		}
	}
}

func TestContainsAll(t *testing.T) {
	tests := []struct {
		value    []int
		required []int
		wantErr  bool
	}{
		{[]int{1, 2, 3}, []int{1, 2}, false},
		{[]int{1, 2, 3}, []int{1, 2, 3}, false},
		{[]int{1, 2}, []int{1, 2, 3}, true},
	}
	for _, tt := range tests {
		v := ContainsAll(tt.value, tt.required, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("ContainsAll(%v, %v) failed = %v, wantErr %v", tt.value, tt.required, v.Failed(), tt.wantErr)
		}
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		value   []int
		options []int
		wantErr bool
	}{
		{[]int{1, 2, 3}, []int{2, 4}, false},
		{[]int{1, 2, 3}, []int{4, 5}, true},
	}
	for _, tt := range tests {
		v := ContainsAny(tt.value, tt.options, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("ContainsAny(%v, %v) failed = %v, wantErr %v", tt.value, tt.options, v.Failed(), tt.wantErr)
		}
	}
}

func TestContainsNone(t *testing.T) {
	tests := []struct {
		value     []int
		forbidden []int
		wantErr   bool
	}{
		{[]int{1, 2, 3}, []int{4, 5}, false},
		{[]int{1, 2, 3}, []int{2, 4}, true},
	}
	for _, tt := range tests {
		v := ContainsNone(tt.value, tt.forbidden, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("ContainsNone(%v, %v) failed = %v, wantErr %v", tt.value, tt.forbidden, v.Failed(), tt.wantErr)
		}
	}
}

func TestEach(t *testing.T) {
	t.Run("all pass", func(t *testing.T) {
		values := []int{2, 4, 6}
		r := Each(values, func(v int, i int) *Validation {
			field := fmt.Sprintf("items[%d]", i)
			if v%2 != 0 {
				return validation(fieldErr(field, "is odd"), field, "even")
			}
			return validation(nil, field, "even")
		})
		if r.Err() != nil {
			t.Errorf("Each expected nil error, got %v", r.Err())
		}
	})

	t.Run("some fail", func(t *testing.T) {
		values := []int{2, 3, 4}
		r := Each(values, func(v int, i int) *Validation {
			field := fmt.Sprintf("items[%d]", i)
			if v%2 != 0 {
				return validation(fieldErr(field, "is odd"), field, "even")
			}
			return validation(nil, field, "even")
		})
		if r.Err() == nil {
			t.Error("Each expected error, got nil")
		}
	})

	t.Run("tracks validators", func(t *testing.T) {
		values := []int{1, 2}
		r := Each(values, func(_ int, i int) *Validation {
			return validation(nil, fmt.Sprintf("items[%d]", i), "positive")
		})
		if !r.HasValidator("items[0]", "positive") {
			t.Error("should track validator for items[0]")
		}
		if !r.HasValidator("items[1]", "positive") {
			t.Error("should track validator for items[1]")
		}
	})
}

func TestAllSatisfy(t *testing.T) {
	t.Run("all satisfy", func(t *testing.T) {
		values := []int{2, 4, 6}
		v := AllSatisfy(values, func(v int) bool { return v%2 == 0 }, "field", "must be even")
		if v.Failed() {
			t.Errorf("AllSatisfy expected pass, got fail")
		}
	})

	t.Run("not all satisfy", func(t *testing.T) {
		values := []int{2, 3, 4}
		v := AllSatisfy(values, func(v int) bool { return v%2 == 0 }, "field", "must be even")
		if !v.Failed() {
			t.Error("AllSatisfy expected fail, got pass")
		}
	})
}

func TestAnySatisfies(t *testing.T) {
	t.Run("one satisfies", func(t *testing.T) {
		values := []int{1, 2, 3}
		v := AnySatisfies(values, func(v int) bool { return v%2 == 0 }, "field", "must have even")
		if v.Failed() {
			t.Errorf("AnySatisfies expected pass, got fail")
		}
	})

	t.Run("none satisfy", func(t *testing.T) {
		values := []int{1, 3, 5}
		v := AnySatisfies(values, func(v int) bool { return v%2 == 0 }, "field", "must have even")
		if !v.Failed() {
			t.Error("AnySatisfies expected fail, got pass")
		}
	})
}

func TestSubset(t *testing.T) {
	tests := []struct {
		value    []int
		superset []int
		wantErr  bool
	}{
		{[]int{1, 2}, []int{1, 2, 3}, false},
		{[]int{1, 2, 3}, []int{1, 2, 3}, false},
		{[]int{}, []int{1, 2, 3}, false},
		{[]int{1, 4}, []int{1, 2, 3}, true},
	}
	for _, tt := range tests {
		v := Subset(tt.value, tt.superset, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Subset(%v, %v) failed = %v, wantErr %v", tt.value, tt.superset, v.Failed(), tt.wantErr)
		}
	}
}

func TestDisjoint(t *testing.T) {
	tests := []struct {
		value   []int
		other   []int
		wantErr bool
	}{
		{[]int{1, 2}, []int{3, 4}, false},
		{[]int{1, 2}, []int{2, 3}, true},
	}
	for _, tt := range tests {
		v := Disjoint(tt.value, tt.other, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Disjoint(%v, %v) failed = %v, wantErr %v", tt.value, tt.other, v.Failed(), tt.wantErr)
		}
	}
}

func TestEachValue(t *testing.T) {
	t.Run("all pass", func(t *testing.T) {
		values := []string{"hello", "world"}
		r := EachValue(values, func(v string) *Validation {
			return MinLen(v, 3, "item")
		})
		if r.Err() != nil {
			t.Errorf("EachValue expected nil error, got %v", r.Err())
		}
	})

	t.Run("some fail", func(t *testing.T) {
		values := []string{"hello", "hi", "world"}
		r := EachValue(values, func(v string) *Validation {
			return MinLen(v, 3, "item")
		})
		if r.Err() == nil {
			t.Error("EachValue expected error, got nil")
		}
	})

	t.Run("tracks validators", func(t *testing.T) {
		values := []string{"a", "b"}
		r := EachValue(values, func(v string) *Validation {
			return Required(v, "item")
		})
		if !r.HasValidator("item", "required") {
			t.Error("should track validator for item")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		var values []string
		r := EachValue(values, func(v string) *Validation {
			return Required(v, "item")
		})
		if r.Err() != nil {
			t.Error("empty slice should have no errors")
		}
	})
}
