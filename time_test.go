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
			v := Before(tt.value, tt.ref, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("Before() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
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
			v := After(tt.value, tt.ref, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("After() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
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
			v := BeforeOrEqual(tt.value, tt.ref, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("BeforeOrEqual() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
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
			v := BetweenTime(tt.value, start, end, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("BetweenTime() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
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
			v := WithinDuration(tt.value, tt.dur, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("WithinDuration() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
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
			v := SameDay(tt.value, tt.ref, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("SameDay() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
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
			v := Weekday(tt.value, tt.day, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("Weekday() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
			}
		})
	}
}

func TestNotWeekend(t *testing.T) {
	monday := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)   // Monday
	saturday := time.Date(2024, 1, 6, 12, 0, 0, 0, time.UTC) // Saturday

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
			v := NotWeekend(tt.value, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("NotWeekend() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
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
			v := NotZeroTime(tt.value, "field")
			if v.Failed() != tt.wantErr {
				t.Errorf("NotZeroTime() failed = %v, wantErr %v", v.Failed(), tt.wantErr)
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
		v := DurationMin(tt.value, tt.min, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("DurationMin(%v, %v) failed = %v, wantErr %v", tt.value, tt.min, v.Failed(), tt.wantErr)
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
		v := DurationMax(tt.value, tt.max, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("DurationMax(%v, %v) failed = %v, wantErr %v", tt.value, tt.max, v.Failed(), tt.wantErr)
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
		v := DurationPositive(tt.value, "field")
		if v.Failed() != tt.wantErr {
			t.Errorf("DurationPositive(%v) failed = %v, wantErr %v", tt.value, v.Failed(), tt.wantErr)
		}
	}
}
