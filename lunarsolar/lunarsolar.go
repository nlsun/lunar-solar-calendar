// Package lunarsolar provides Lunar and Solar date conversions.
//
// A good source of date conversions to verify against is
// https://www.hko.gov.hk/en/gts/time/conversion.htm
package lunarsolar

import (
	"time"

	"github.com/isee15/Lunar-Solar-Calendar-Converter/Go/lunarsolar"
)

type LunarTime struct {
	// Solar calendar time
	time time.Time
	// If true, it's a lunar leap year, and this month is being repeated.
	isLeap bool
}

const (
	day = time.Hour * 24
)

func NewLunarTime(t time.Time, isLeap bool) LunarTime {
	return LunarTime{
		time:   t,
		isLeap: isLeap,
	}
}

func SolarToLunar(t time.Time) LunarTime {
	year, month, day := t.Date()
	lunar := lunarsolar.SolarToLunar(lunarsolar.Solar{
		SolarYear:  year,
		SolarMonth: int(month),
		SolarDay:   day,
	})
	return LunarTime{
		time: time.Date(
			lunar.LunarYear,
			time.Month(lunar.LunarMonth),
			lunar.LunarDay,
			t.Hour(),
			t.Minute(),
			t.Second(),
			t.Nanosecond(),
			t.Location()),
		isLeap: lunar.IsLeap,
	}
}

func LunarToSolar(t LunarTime) time.Time {
	year, month, day := t.time.Date()
	solar := lunarsolar.LunarToSolar(lunarsolar.Lunar{
		IsLeap:     t.isLeap,
		LunarYear:  year,
		LunarMonth: int(month),
		LunarDay:   day,
	})
	return time.Date(solar.SolarYear,
		time.Month(solar.SolarMonth),
		solar.SolarDay,
		t.time.Hour(),
		t.time.Minute(),
		t.time.Second(),
		t.time.Nanosecond(),
		t.time.Location())
}

// IsLunarLeapMonthPossible takes a lunar date that lacks leap year information
// and tries to figure out if it was possible that month was repeated.
func IsLunarLeapMonthPossible(t LunarTime) bool {
	// Take the lunar time assuming it's not a leap month
	t.isLeap = false

	// Assuming a month can't be longer than 31 days, we jump forward by that
	// amount plus 1 to land in the next month.
	diff := 31 - t.time.Day()
	lunarPlus := t.Add(day * time.Duration(diff+1))

	return lunarPlus.isLeap
}

// Solar calendar time
func (t LunarTime) Time() time.Time {
	return t.time
}

// If true, it's a lunar leap year, and this month is being repeated.
func (t LunarTime) IsLeap() bool {
	return t.isLeap
}

func (t LunarTime) Add(d time.Duration) LunarTime {
	return SolarToLunar(LunarToSolar(t).Add(d))
}

func (t LunarTime) AddDate(years int, months int, days int) LunarTime {
	return LunarTime{
		time:   t.time.AddDate(years, months, days),
		isLeap: t.isLeap,
	}
}

func (t LunarTime) Sub(u LunarTime) time.Duration {
	return LunarToSolar(t).Sub(LunarToSolar(u))
}

func (t LunarTime) Equal(u LunarTime) bool {
	return t.time.Equal(u.time) && t.isLeap == u.isLeap
}

func (t LunarTime) Before(u LunarTime) bool {
	return LunarToSolar(t).Before(LunarToSolar(u))
}

func (t LunarTime) After(u LunarTime) bool {
	return LunarToSolar(t).After(LunarToSolar(u))
}

func (t LunarTime) AsLeap(isLeap bool) LunarTime {
	return LunarTime{
		time:   t.time,
		isLeap: isLeap,
	}
}
