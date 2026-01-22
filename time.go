package check

import (
	"time"
)

// Before validates that a time is before the given time.
func Before(v, t time.Time, field string) error {
	if !v.Before(t) {
		return fieldErrf(field, "must be before %s", t.Format(time.RFC3339))
	}
	return nil
}

// After validates that a time is after the given time.
func After(v, t time.Time, field string) error {
	if !v.After(t) {
		return fieldErrf(field, "must be after %s", t.Format(time.RFC3339))
	}
	return nil
}

// BeforeOrEqual validates that a time is before or equal to the given time.
func BeforeOrEqual(v, t time.Time, field string) error {
	if v.After(t) {
		return fieldErrf(field, "must be before or equal to %s", t.Format(time.RFC3339))
	}
	return nil
}

// AfterOrEqual validates that a time is after or equal to the given time.
func AfterOrEqual(v, t time.Time, field string) error {
	if v.Before(t) {
		return fieldErrf(field, "must be after or equal to %s", t.Format(time.RFC3339))
	}
	return nil
}

// BeforeNow validates that a time is before the current time.
func BeforeNow(v time.Time, field string) error {
	if !v.Before(time.Now()) {
		return fieldErr(field, "must be in the past")
	}
	return nil
}

// AfterNow validates that a time is after the current time.
func AfterNow(v time.Time, field string) error {
	if !v.After(time.Now()) {
		return fieldErr(field, "must be in the future")
	}
	return nil
}

// BeforeOrEqualNow validates that a time is before or equal to the current time.
func BeforeOrEqualNow(v time.Time, field string) error {
	if v.After(time.Now()) {
		return fieldErr(field, "must not be in the future")
	}
	return nil
}

// AfterOrEqualNow validates that a time is after or equal to the current time.
func AfterOrEqualNow(v time.Time, field string) error {
	if v.Before(time.Now()) {
		return fieldErr(field, "must not be in the past")
	}
	return nil
}

// InPast is an alias for BeforeNow.
func InPast(v time.Time, field string) error {
	return BeforeNow(v, field)
}

// InFuture is an alias for AfterNow.
func InFuture(v time.Time, field string) error {
	return AfterNow(v, field)
}

// BetweenTime validates that a time is within a range (inclusive).
func BetweenTime(v, start, end time.Time, field string) error {
	if v.Before(start) || v.After(end) {
		return fieldErrf(field, "must be between %s and %s", start.Format(time.RFC3339), end.Format(time.RFC3339))
	}
	return nil
}

// BetweenTimeExclusive validates that a time is within a range (exclusive).
func BetweenTimeExclusive(v, start, end time.Time, field string) error {
	if !v.After(start) || !v.Before(end) {
		return fieldErrf(field, "must be between %s and %s (exclusive)", start.Format(time.RFC3339), end.Format(time.RFC3339))
	}
	return nil
}

// WithinDuration validates that a time is within a duration from now.
func WithinDuration(v time.Time, d time.Duration, field string) error {
	now := time.Now()
	diff := v.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	if diff > d {
		return fieldErrf(field, "must be within %s of now", d)
	}
	return nil
}

// WithinDurationOf validates that a time is within a duration of a reference time.
func WithinDurationOf(v time.Time, d time.Duration, ref time.Time, field string) error {
	diff := v.Sub(ref)
	if diff < 0 {
		diff = -diff
	}
	if diff > d {
		return fieldErrf(field, "must be within %s of reference time", d)
	}
	return nil
}

// SameDay validates that a time is on the same day as the reference time.
func SameDay(v, ref time.Time, field string) error {
	vYear, vMonth, vDay := v.Date()
	refYear, refMonth, refDay := ref.Date()
	if vYear != refYear || vMonth != refMonth || vDay != refDay {
		return fieldErr(field, "must be on the same day")
	}
	return nil
}

// SameMonth validates that a time is in the same month as the reference time.
func SameMonth(v, ref time.Time, field string) error {
	vYear, vMonth, _ := v.Date()
	refYear, refMonth, _ := ref.Date()
	if vYear != refYear || vMonth != refMonth {
		return fieldErr(field, "must be in the same month")
	}
	return nil
}

// SameYear validates that a time is in the same year as the reference time.
func SameYear(v, ref time.Time, field string) error {
	if v.Year() != ref.Year() {
		return fieldErr(field, "must be in the same year")
	}
	return nil
}

// Weekday validates that a time is on the specified weekday.
func Weekday(v time.Time, day time.Weekday, field string) error {
	if v.Weekday() != day {
		return fieldErrf(field, "must be on a %s", day)
	}
	return nil
}

// WeekdayIn validates that a time is on one of the specified weekdays.
func WeekdayIn(v time.Time, days []time.Weekday, field string) error {
	vDay := v.Weekday()
	for _, day := range days {
		if vDay == day {
			return nil
		}
	}
	return fieldErr(field, "must be on an allowed weekday")
}

// NotWeekend validates that a time is not on Saturday or Sunday.
func NotWeekend(v time.Time, field string) error {
	day := v.Weekday()
	if day == time.Saturday || day == time.Sunday {
		return fieldErr(field, "must not be on a weekend")
	}
	return nil
}

// IsWeekend validates that a time is on Saturday or Sunday.
func IsWeekend(v time.Time, field string) error {
	day := v.Weekday()
	if day != time.Saturday && day != time.Sunday {
		return fieldErr(field, "must be on a weekend")
	}
	return nil
}

// NotZeroTime validates that a time is not the zero value.
func NotZeroTime(v time.Time, field string) error {
	if v.IsZero() {
		return fieldErr(field, "must not be empty")
	}
	return nil
}

// ZeroTime validates that a time is the zero value.
func ZeroTime(v time.Time, field string) error {
	if !v.IsZero() {
		return fieldErr(field, "must be empty")
	}
	return nil
}

// TimeInTimezone validates that a time's location matches the expected timezone.
func TimeInTimezone(v time.Time, loc *time.Location, field string) error {
	if v.Location().String() != loc.String() {
		return fieldErrf(field, "must be in timezone %s", loc)
	}
	return nil
}

// DurationMin validates that a duration is at least the minimum.
func DurationMin(v, minDur time.Duration, field string) error {
	if v < minDur {
		return fieldErrf(field, "must be at least %s", minDur)
	}
	return nil
}

// DurationMax validates that a duration is at most the maximum.
func DurationMax(v, maxDur time.Duration, field string) error {
	if v > maxDur {
		return fieldErrf(field, "must be at most %s", maxDur)
	}
	return nil
}

// DurationBetween validates that a duration is within a range (inclusive).
func DurationBetween(v, minDur, maxDur time.Duration, field string) error {
	if v < minDur || v > maxDur {
		return fieldErrf(field, "must be between %s and %s", minDur, maxDur)
	}
	return nil
}

// DurationPositive validates that a duration is positive.
func DurationPositive(v time.Duration, field string) error {
	if v <= 0 {
		return fieldErr(field, "must be positive")
	}
	return nil
}

// DurationNonNegative validates that a duration is non-negative.
func DurationNonNegative(v time.Duration, field string) error {
	if v < 0 {
		return fieldErr(field, "must not be negative")
	}
	return nil
}
