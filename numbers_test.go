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
		err := Min(tt.value, tt.min, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Min(%d, %d) = %v, wantErr %v", tt.value, tt.min, err, tt.wantErr)
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
		err := Max(tt.value, tt.max, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Max(%d, %d) = %v, wantErr %v", tt.value, tt.max, err, tt.wantErr)
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
		err := Between(tt.value, tt.min, tt.max, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Between(%d, %d, %d) = %v, wantErr %v", tt.value, tt.min, tt.max, err, tt.wantErr)
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
		err := BetweenExclusive(tt.value, tt.min, tt.max, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("BetweenExclusive(%d, %d, %d) = %v, wantErr %v", tt.value, tt.min, tt.max, err, tt.wantErr)
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
		err := Positive(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Positive(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := Negative(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Negative(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := NonNegative(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NonNegative(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := NonPositive(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NonPositive(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := Zero(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Zero(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := NonZero(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("NonZero(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := MultipleOf(tt.value, tt.divisor, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("MultipleOf(%d, %d) = %v, wantErr %v", tt.value, tt.divisor, err, tt.wantErr)
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
		err := Even(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Even(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := Odd(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Odd(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := OneOfValues(tt.value, allowed, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("OneOfValues(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := GreaterThan(tt.value, tt.threshold, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("GreaterThan(%d, %d) = %v, wantErr %v", tt.value, tt.threshold, err, tt.wantErr)
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
		err := LessThan(tt.value, tt.threshold, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("LessThan(%d, %d) = %v, wantErr %v", tt.value, tt.threshold, err, tt.wantErr)
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
		err := Percentage(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Percentage(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := PortNumber(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("PortNumber(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := HTTPStatusCode(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("HTTPStatusCode(%d) = %v, wantErr %v", tt.value, err, tt.wantErr)
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
		err := Min(tt.value, tt.min, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("Min(%f, %f) = %v, wantErr %v", tt.value, tt.min, err, tt.wantErr)
		}
	}
}
