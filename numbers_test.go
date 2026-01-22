package check

import (
	"testing"
)

func TestMin(t *testing.T) {
	tests := []struct {
		value   int
		min     int
		wantErr bool
	}{
		{10, 5, false},
		{5, 5, false},
		{4, 5, true},
	}
	for _, tt := range tests {
		v := Min(tt.value, tt.min, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Min(%d, %d) failed = %v, wantErr %v", tt.value, tt.min, v.Failed(), tt.wantErr)
		}
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		value   int
		max     int
		wantErr bool
	}{
		{5, 10, false},
		{10, 10, false},
		{11, 10, true},
	}
	for _, tt := range tests {
		v := Max(tt.value, tt.max, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Max(%d, %d) failed = %v, wantErr %v", tt.value, tt.max, v.Failed(), tt.wantErr)
		}
	}
}

func TestBetween(t *testing.T) {
	tests := []struct {
		value   int
		min     int
		max     int
		wantErr bool
	}{
		{5, 1, 10, false},
		{1, 1, 10, false},
		{10, 1, 10, false},
		{0, 1, 10, true},
		{11, 1, 10, true},
	}
	for _, tt := range tests {
		v := Between(tt.value, tt.min, tt.max, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Between(%d, %d, %d) failed = %v, wantErr %v", tt.value, tt.min, tt.max, v.Failed(), tt.wantErr)
		}
	}
}

func TestBetweenExclusive(t *testing.T) {
	tests := []struct {
		value   int
		min     int
		max     int
		wantErr bool
	}{
		{5, 1, 10, false},
		{1, 1, 10, true},  // min excluded
		{10, 1, 10, true}, // max excluded
		{2, 1, 10, false},
		{9, 1, 10, false},
	}
	for _, tt := range tests {
		v := BetweenExclusive(tt.value, tt.min, tt.max, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("BetweenExclusive(%d, %d, %d) failed = %v, wantErr %v", tt.value, tt.min, tt.max, v.Failed(), tt.wantErr)
		}
	}
}

func TestPositive(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{1, false},
		{100, false},
		{0, true},
		{-1, true},
	}
	for _, tt := range tests {
		v := Positive(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Positive(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestNegative(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{-1, false},
		{-100, false},
		{0, true},
		{1, true},
	}
	for _, tt := range tests {
		v := Negative(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Negative(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestNonNegative(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{0, false},
		{1, false},
		{-1, true},
	}
	for _, tt := range tests {
		v := NonNegative(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NonNegative(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestNonPositive(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{0, false},
		{-1, false},
		{1, true},
	}
	for _, tt := range tests {
		v := NonPositive(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NonPositive(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestZero(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{0, false},
		{1, true},
		{-1, true},
	}
	for _, tt := range tests {
		v := Zero(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Zero(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestNonZero(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{1, false},
		{-1, false},
		{0, true},
	}
	for _, tt := range tests {
		v := NonZero(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("NonZero(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestMultipleOf(t *testing.T) {
	tests := []struct {
		value   int
		divisor int
		wantErr bool
	}{
		{10, 5, false},
		{10, 2, false},
		{10, 3, true},
	}
	for _, tt := range tests {
		v := MultipleOf(tt.value, tt.divisor, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("MultipleOf(%d, %d) failed = %v, wantErr %v", tt.value, tt.divisor, v.Failed(), tt.wantErr)
		}
	}
}

func TestEven(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{0, false},
		{2, false},
		{-4, false},
		{1, true},
		{-3, true},
	}
	for _, tt := range tests {
		v := Even(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Even(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestOdd(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{1, false},
		{-3, false},
		{0, true},
		{2, true},
	}
	for _, tt := range tests {
		v := Odd(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Odd(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestOneOfValues(t *testing.T) {
	allowed := []int{1, 2, 3}
	tests := []struct {
		value   int
		wantErr bool
	}{
		{1, false},
		{2, false},
		{3, false},
		{4, true},
	}
	for _, tt := range tests {
		v := OneOfValues(tt.value, allowed, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("OneOfValues(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	tests := []struct {
		value     int
		threshold int
		wantErr   bool
	}{
		{6, 5, false},
		{5, 5, true},
		{4, 5, true},
	}
	for _, tt := range tests {
		v := GreaterThan(tt.value, tt.threshold, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("GreaterThan(%d, %d) failed = %v, wantErr %v", tt.value, tt.threshold, v.Failed(), tt.wantErr)
		}
	}
}

func TestLessThan(t *testing.T) {
	tests := []struct {
		value     int
		threshold int
		wantErr   bool
	}{
		{4, 5, false},
		{5, 5, true},
		{6, 5, true},
	}
	for _, tt := range tests {
		v := LessThan(tt.value, tt.threshold, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("LessThan(%d, %d) failed = %v, wantErr %v", tt.value, tt.threshold, v.Failed(), tt.wantErr)
		}
	}
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{0, false},
		{50, false},
		{100, false},
		{-1, true},
		{101, true},
	}
	for _, tt := range tests {
		v := Percentage(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Percentage(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestPortNumber(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{1, false},
		{80, false},
		{65535, false},
		{0, true},
		{65536, true},
	}
	for _, tt := range tests {
		v := PortNumber(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("PortNumber(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestHTTPStatusCode(t *testing.T) {
	tests := []struct {
		value   int
		wantErr bool
	}{
		{100, false},
		{200, false},
		{404, false},
		{500, false},
		{599, false},
		{99, true},
		{600, true},
	}
	for _, tt := range tests {
		v := HTTPStatusCode(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("HTTPStatusCode(%d) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}

func TestMinFloat(t *testing.T) {
	tests := []struct {
		value   float64
		min     float64
		wantErr bool
	}{
		{5.5, 5.0, false},
		{5.0, 5.0, false},
		{4.9, 5.0, true},
	}
	for _, tt := range tests {
		v := Min(tt.value, tt.min, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("Min(%f, %f) failed = %v, wantErr %v", tt.value, tt.min, v.Failed(), tt.wantErr)
		}
	}
}
