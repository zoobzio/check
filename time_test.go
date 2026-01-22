package check

import (
	"testing"
	"time"
)

func TestBefore(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name    string
		value   time.Time
		ref     time.Time
		wantErr bool
	}{
		{"past before now", past, now, false},
		{"now before future", now, future, false},
		{"future before now", future, now, true},
		{"same time", now, now, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Before(tt.value, tt.ref, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("Before() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAfter(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name    string
		value   time.Time
		ref     time.Time
		wantErr bool
	}{
		{"future after now", future, now, false},
		{"now after past", now, past, false},
		{"past after now", past, now, true},
		{"same time", now, now, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := After(tt.value, tt.ref, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("After() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBeforeOrEqual(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name    string
		value   time.Time
		ref     time.Time
		wantErr bool
	}{
		{"past before now", past, now, false},
		{"same time", now, now, false},
		{"future before now", future, now, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BeforeOrEqual(tt.value, tt.ref, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("BeforeOrEqual() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBetweenTime(t *testing.T) {
	now := time.Now()
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	before := now.Add(-2 * time.Hour)
	after := now.Add(2 * time.Hour)

	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"in range", now, false},
		{"at start", start, false},
		{"at end", end, false},
		{"before range", before, true},
		{"after range", after, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BetweenTime(tt.value, start, end, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("BetweenTime() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithinDuration(t *testing.T) {
	now := time.Now()
	nearby := now.Add(30 * time.Minute)
	far := now.Add(2 * time.Hour)

	tests := []struct {
		name    string
		value   time.Time
		dur     time.Duration
		wantErr bool
	}{
		{"within duration", nearby, time.Hour, false},
		{"at boundary", now.Add(time.Hour), time.Hour, false},
		{"beyond duration", far, time.Hour, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WithinDuration(tt.value, tt.dur, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("WithinDuration() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSameDay(t *testing.T) {
	now := time.Now()
	sameDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	nextDay := now.AddDate(0, 0, 1)

	tests := []struct {
		name    string
		value   time.Time
		ref     time.Time
		wantErr bool
	}{
		{"same day", now, sameDay, false},
		{"different day", now, nextDay, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SameDay(tt.value, tt.ref, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("SameDay() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWeekday(t *testing.T) {
	// Find a known Monday
	monday := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC) // Jan 1, 2024 is Monday

	tests := []struct {
		name    string
		value   time.Time
		day     time.Weekday
		wantErr bool
	}{
		{"correct weekday", monday, time.Monday, false},
		{"wrong weekday", monday, time.Tuesday, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Weekday(tt.value, tt.day, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("Weekday() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotWeekend(t *testing.T) {
	monday := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)    // Monday
	saturday := time.Date(2024, 1, 6, 12, 0, 0, 0, time.UTC)  // Saturday

	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"weekday", monday, false},
		{"weekend", saturday, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NotWeekend(tt.value, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("NotWeekend() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotZeroTime(t *testing.T) {
	tests := []struct {
		name    string
		value   time.Time
		wantErr bool
	}{
		{"not zero", time.Now(), false},
		{"zero", time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NotZeroTime(tt.value, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("NotZeroTime() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDurationMin(t *testing.T) {
	tests := []struct {
		value   time.Duration
		min     time.Duration
		wantErr bool
	}{
		{time.Hour, 30 * time.Minute, false},
		{30 * time.Minute, 30 * time.Minute, false},
		{15 * time.Minute, 30 * time.Minute, true},
	}
	for _, tt := range tests {
		err := DurationMin(tt.value, tt.min, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("DurationMin(%v, %v) = %v, wantErr %v", tt.value, tt.min, err, tt.wantErr)
		}
	}
}

func TestDurationMax(t *testing.T) {
	tests := []struct {
		value   time.Duration
		max     time.Duration
		wantErr bool
	}{
		{30 * time.Minute, time.Hour, false},
		{time.Hour, time.Hour, false},
		{2 * time.Hour, time.Hour, true},
	}
	for _, tt := range tests {
		err := DurationMax(tt.value, tt.max, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("DurationMax(%v, %v) = %v, wantErr %v", tt.value, tt.max, err, tt.wantErr)
		}
	}
}

func TestDurationPositive(t *testing.T) {
	tests := []struct {
		value   time.Duration
		wantErr bool
	}{
		{time.Second, false},
		{0, true},
		{-time.Second, true},
	}
	for _, tt := range tests {
		err := DurationPositive(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("DurationPositive(%v) = %v, wantErr %v", tt.value, err, tt.wantErr)
		}
	}
}
