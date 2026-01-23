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
		{name: "past before now", value: past, ref: now, wantErr: false},
		{name: "now before future", value: now, ref: future, wantErr: false},
		{name: "future before now", value: future, ref: now, wantErr: true},
		{name: "same time", value: now, ref: now, wantErr: true},
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
		{name: "future after now", value: future, ref: now, wantErr: false},
		{name: "now after past", value: now, ref: past, wantErr: false},
		{name: "past after now", value: past, ref: now, wantErr: true},
		{name: "same time", value: now, ref: now, wantErr: true},
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
		{name: "past before now", value: past, ref: now, wantErr: false},
		{name: "same time", value: now, ref: now, wantErr: false},
		{name: "future before now", value: future, ref: now, wantErr: true},
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
		{name: "in range", value: now, wantErr: false},
		{name: "at start", value: start, wantErr: false},
		{name: "at end", value: end, wantErr: false},
		{name: "before range", value: before, wantErr: true},
		{name: "after range", value: after, wantErr: true},
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
		value   time.Time
		name    string
		dur     time.Duration
		wantErr bool
	}{
		{name: "within duration", value: nearby, dur: time.Hour, wantErr: false},
		{name: "at boundary", value: now.Add(time.Hour), dur: time.Hour, wantErr: false},
		{name: "beyond duration", value: far, dur: time.Hour, wantErr: true},
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
		value   time.Time
		ref     time.Time
		name    string
		wantErr bool
	}{
		{name: "same day", value: now, ref: sameDay, wantErr: false},
		{name: "different day", value: now, ref: nextDay, wantErr: true},
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
		value   time.Time
		name    string
		day     time.Weekday
		wantErr bool
	}{
		{name: "correct weekday", value: monday, day: time.Monday, wantErr: false},
		{name: "wrong weekday", value: monday, day: time.Tuesday, wantErr: true},
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
		value   time.Time
		name    string
		wantErr bool
	}{
		{name: "weekday", value: monday, wantErr: false},
		{name: "weekend", value: saturday, wantErr: true},
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
		value   time.Time
		name    string
		wantErr bool
	}{
		{name: "not zero", value: time.Now(), wantErr: false},
		{name: "zero", value: time.Time{}, wantErr: true},
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
