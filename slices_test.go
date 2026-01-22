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
		err := NotEmpty(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NotEmpty(%v) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := Empty(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Empty(%v) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := MinItems(tt.value, tt.min, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("MinItems(%v, %d) = %v, wantErr %v", tt.value, tt.min, err, tt.wantErr)
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
		err := MaxItems(tt.value, tt.max, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("MaxItems(%v, %d) = %v, wantErr %v", tt.value, tt.max, err, tt.wantErr)
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
		err := ExactItems(tt.value, tt.count, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("ExactItems(%v, %d) = %v, wantErr %v", tt.value, tt.count, err, tt.wantErr)
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
		err := Unique(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Unique(%v) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := SliceContains(tt.value, tt.elem, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("SliceContains(%v, %d) = %v, wantErr %v", tt.value, tt.elem, err, tt.wantErr)
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
		err := SliceNotContains(tt.value, tt.elem, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("SliceNotContains(%v, %d) = %v, wantErr %v", tt.value, tt.elem, err, tt.wantErr)
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
		err := ContainsAll(tt.value, tt.required, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("ContainsAll(%v, %v) = %v, wantErr %v", tt.value, tt.required, err, tt.wantErr)
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
		err := ContainsAny(tt.value, tt.options, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("ContainsAny(%v, %v) = %v, wantErr %v", tt.value, tt.options, err, tt.wantErr)
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
		err := ContainsNone(tt.value, tt.forbidden, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("ContainsNone(%v, %v) = %v, wantErr %v", tt.value, tt.forbidden, err, tt.wantErr)
		}
	}
}

func TestEach(t *testing.T) {
	t.Run("all pass", func(t *testing.T) {
		values := []int{2, 4, 6}
		err := Each(values, func(v int, i int) error {
			if v%2 != 0 {
				return fmt.Errorf("item %d is odd", i)
			}
			return nil
		})
		if err != nil {
			t.Errorf("Each expected nil, got %v", err)
		}
	})

	t.Run("some fail", func(t *testing.T) {
		values := []int{2, 3, 4}
		err := Each(values, func(v int, i int) error {
			if v%2 != 0 {
				return fmt.Errorf("item %d is odd", i)
			}
			return nil
		})
		if err == nil {
			t.Error("Each expected error, got nil")
		}
	})
}

func TestAllSatisfy(t *testing.T) {
	t.Run("all satisfy", func(t *testing.T) {
		values := []int{2, 4, 6}
		err := AllSatisfy(values, func(v int) bool { return v%2 == 0 }, "field", "must be even")
		if err != nil {
			t.Errorf("AllSatisfy expected nil, got %v", err)
		}
	})

	t.Run("not all satisfy", func(t *testing.T) {
		values := []int{2, 3, 4}
		err := AllSatisfy(values, func(v int) bool { return v%2 == 0 }, "field", "must be even")
		if err == nil {
			t.Error("AllSatisfy expected error, got nil")
		}
	})
}

func TestAnySatisfies(t *testing.T) {
	t.Run("one satisfies", func(t *testing.T) {
		values := []int{1, 2, 3}
		err := AnySatisfies(values, func(v int) bool { return v%2 == 0 }, "field", "must have even")
		if err != nil {
			t.Errorf("AnySatisfies expected nil, got %v", err)
		}
	})

	t.Run("none satisfy", func(t *testing.T) {
		values := []int{1, 3, 5}
		err := AnySatisfies(values, func(v int) bool { return v%2 == 0 }, "field", "must have even")
		if err == nil {
			t.Error("AnySatisfies expected error, got nil")
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
		err := Subset(tt.value, tt.superset, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Subset(%v, %v) = %v, wantErr %v", tt.value, tt.superset, err, tt.wantErr)
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
		err := Disjoint(tt.value, tt.other, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Disjoint(%v, %v) = %v, wantErr %v", tt.value, tt.other, err, tt.wantErr)
		}
	}
}
