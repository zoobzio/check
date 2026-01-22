package check

import (
	"time"
)

// Before validates that a time is before the given time.
func Before(v, t time.Time, field string) *Validation {
	var err error
	if !v.Before(t) {
		err = fieldErrf(field, "must be before %s", t.Format(time.RFC3339))
	}
	return validation(err, field, "before")
}

// After validates that a time is after the given time.
func After(v, t time.Time, field string) *Validation {
	var err error
	if !v.After(t) {
		err = fieldErrf(field, "must be after %s", t.Format(time.RFC3339))
	}
	return validation(err, field, "after")
}

// BeforeOrEqual validates that a time is before or equal to the given time.
func BeforeOrEqual(v, t time.Time, field string) *Validation {
	var err error
	if v.After(t) {
		err = fieldErrf(field, "must be before or equal to %s", t.Format(time.RFC3339))
	}
	return validation(err, field, "lte")
}

// AfterOrEqual validates that a time is after or equal to the given time.
func AfterOrEqual(v, t time.Time, field string) *Validation {
	var err error
	if v.Before(t) {
		err = fieldErrf(field, "must be after or equal to %s", t.Format(time.RFC3339))
	}
	return validation(err, field, "gte")
}

// BeforeNow validates that a time is before the current time.
func BeforeNow(v time.Time, field string) *Validation {
	var err error
	if !v.Before(time.Now()) {
		err = fieldErr(field, "must be in the past")
	}
	return validation(err, field, "past")
}

// AfterNow validates that a time is after the current time.
func AfterNow(v time.Time, field string) *Validation {
	var err error
	if !v.After(time.Now()) {
		err = fieldErr(field, "must be in the future")
	}
	return validation(err, field, "future")
}

// BeforeOrEqualNow validates that a time is before or equal to the current time.
func BeforeOrEqualNow(v time.Time, field string) *Validation {
	var err error
	if v.After(time.Now()) {
		err = fieldErr(field, "must not be in the future")
	}
	return validation(err, field, "pastoreq")
}

// AfterOrEqualNow validates that a time is after or equal to the current time.
func AfterOrEqualNow(v time.Time, field string) *Validation {
	var err error
	if v.Before(time.Now()) {
		err = fieldErr(field, "must not be in the past")
	}
	return validation(err, field, "futureoreq")
}

// InPast is an alias for BeforeNow.
func InPast(v time.Time, field string) *Validation {
	return BeforeNow(v, field)
}

// InFuture is an alias for AfterNow.
func InFuture(v time.Time, field string) *Validation {
	return AfterNow(v, field)
}

// BetweenTime validates that a time is within a range (inclusive).
func BetweenTime(v, start, end time.Time, field string) *Validation {
	var err error
	if v.Before(start) || v.After(end) {
		err = fieldErrf(field, "must be between %s and %s", start.Format(time.RFC3339), end.Format(time.RFC3339))
	}
	return validation(err, field, "after", "before")
}

// BetweenTimeExclusive validates that a time is within a range (exclusive).
func BetweenTimeExclusive(v, start, end time.Time, field string) *Validation {
	var err error
	if !v.After(start) || !v.Before(end) {
		err = fieldErrf(field, "must be between %s and %s (exclusive)", start.Format(time.RFC3339), end.Format(time.RFC3339))
	}
	return validation(err, field, "gt", "lt")
}

// WithinDuration validates that a time is within a duration from now.
func WithinDuration(v time.Time, d time.Duration, field string) *Validation {
	var err error
	now := time.Now()
	diff := v.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	if diff > d {
		err = fieldErrf(field, "must be within %s of now", d)
	}
	return validation(err, field, "within")
}

// WithinDurationOf validates that a time is within a duration of a reference time.
func WithinDurationOf(v time.Time, d time.Duration, ref time.Time, field string) *Validation {
	var err error
	diff := v.Sub(ref)
	if diff < 0 {
		diff = -diff
	}
	if diff > d {
		err = fieldErrf(field, "must be within %s of reference time", d)
	}
	return validation(err, field, "within")
}

// SameDay validates that a time is on the same day as the reference time.
func SameDay(v, ref time.Time, field string) *Validation {
	var err error
	vYear, vMonth, vDay := v.Date()
	refYear, refMonth, refDay := ref.Date()
	if vYear != refYear || vMonth != refMonth || vDay != refDay {
		err = fieldErr(field, "must be on the same day")
	}
	return validation(err, field, "sameday")
}

// SameMonth validates that a time is in the same month as the reference time.
func SameMonth(v, ref time.Time, field string) *Validation {
	var err error
	vYear, vMonth, _ := v.Date()
	refYear, refMonth, _ := ref.Date()
	if vYear != refYear || vMonth != refMonth {
		err = fieldErr(field, "must be in the same month")
	}
	return validation(err, field, "samemonth")
}

// SameYear validates that a time is in the same year as the reference time.
func SameYear(v, ref time.Time, field string) *Validation {
	var err error
	if v.Year() != ref.Year() {
		err = fieldErr(field, "must be in the same year")
	}
	return validation(err, field, "sameyear")
}

// Weekday validates that a time is on the specified weekday.
func Weekday(v time.Time, day time.Weekday, field string) *Validation {
	var err error
	if v.Weekday() != day {
		err = fieldErrf(field, "must be on a %s", day)
	}
	return validation(err, field, "weekday")
}

// WeekdayIn validates that a time is on one of the specified weekdays.
func WeekdayIn(v time.Time, days []time.Weekday, field string) *Validation {
	var err error
	vDay := v.Weekday()
	found := false
	for _, day := range days {
		if vDay == day {
			found = true
			break
		}
	}
	if !found {
		err = fieldErr(field, "must be on an allowed weekday")
	}
	return validation(err, field, "weekday")
}

// NotWeekend validates that a time is not on Saturday or Sunday.
func NotWeekend(v time.Time, field string) *Validation {
	var err error
	day := v.Weekday()
	if day == time.Saturday || day == time.Sunday {
		err = fieldErr(field, "must not be on a weekend")
	}
	return validation(err, field, "notweekend")
}

// IsWeekend validates that a time is on Saturday or Sunday.
func IsWeekend(v time.Time, field string) *Validation {
	var err error
	day := v.Weekday()
	if day != time.Saturday && day != time.Sunday {
		err = fieldErr(field, "must be on a weekend")
	}
	return validation(err, field, "weekend")
}

// NotZeroTime validates that a time is not the zero value.
func NotZeroTime(v time.Time, field string) *Validation {
	var err error
	if v.IsZero() {
		err = fieldErr(field, "must not be empty")
	}
	return validation(err, field, "required")
}

// ZeroTime validates that a time is the zero value.
func ZeroTime(v time.Time, field string) *Validation {
	var err error
	if !v.IsZero() {
		err = fieldErr(field, "must be empty")
	}
	return validation(err, field, "empty")
}

// TimeInTimezone validates that a time's location matches the expected timezone.
func TimeInTimezone(v time.Time, loc *time.Location, field string) *Validation {
	var err error
	if loc == nil {
		err = fieldErr(field, "timezone must be provided")
	} else if v.Location().String() != loc.String() {
		err = fieldErrf(field, "must be in timezone %s", loc)
	}
	return validation(err, field, "timezone")
}

// DurationMin validates that a duration is at least the minimum.
func DurationMin(v, minDur time.Duration, field string) *Validation {
	var err error
	if v < minDur {
		err = fieldErrf(field, "must be at least %s", minDur)
	}
	return validation(err, field, "min")
}

// DurationMax validates that a duration is at most the maximum.
func DurationMax(v, maxDur time.Duration, field string) *Validation {
	var err error
	if v > maxDur {
		err = fieldErrf(field, "must be at most %s", maxDur)
	}
	return validation(err, field, "max")
}

// DurationBetween validates that a duration is within a range (inclusive).
func DurationBetween(v, minDur, maxDur time.Duration, field string) *Validation {
	var err error
	if v < minDur || v > maxDur {
		err = fieldErrf(field, "must be between %s and %s", minDur, maxDur)
	}
	return validation(err, field, "min", "max")
}

// DurationPositive validates that a duration is positive.
func DurationPositive(v time.Duration, field string) *Validation {
	var err error
	if v <= 0 {
		err = fieldErr(field, "must be positive")
	}
	return validation(err, field, "gt")
}

// DurationNonNegative validates that a duration is non-negative.
func DurationNonNegative(v time.Duration, field string) *Validation {
	var err error
	if v < 0 {
		err = fieldErr(field, "must not be negative")
	}
	return validation(err, field, "gte")
}
