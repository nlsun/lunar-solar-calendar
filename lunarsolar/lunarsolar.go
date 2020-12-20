package lunarsolar

import (
	"time"

	"github.com/isee15/Lunar-Solar-Calendar-Converter/Go/lunarsolar"
)

type LunarTime struct {
	// Solar calendar time
	Time time.Time
	// If true, it's a lunar leap year, and this month is being repeated.
	IsLeap bool
}

const (
	day = time.Hour * 24
)

func SolarToLunar(t time.Time) LunarTime {
	year, month, day := t.Date()
	lunar := lunarsolar.SolarToLunar(lunarsolar.Solar{
		SolarYear:  year,
		SolarMonth: int(month),
		SolarDay:   day,
	})
	return LunarTime{
		Time: time.Date(
			lunar.LunarYear,
			time.Month(lunar.LunarMonth),
			lunar.LunarDay,
			t.Hour(),
			t.Minute(),
			t.Second(),
			t.Nanosecond(),
			t.Location()),
		IsLeap: lunar.IsLeap,
	}
}

func LunarToSolar(t LunarTime) time.Time {
	year, month, day := t.Time.Date()
	solar := lunarsolar.LunarToSolar(lunarsolar.Lunar{
		LunarYear:  year,
		LunarMonth: int(month),
		LunarDay:   day,
	})
	return time.Date(solar.SolarYear,
		time.Month(solar.SolarMonth),
		solar.SolarDay,
		t.Time.Hour(),
		t.Time.Minute(),
		t.Time.Second(),
		t.Time.Nanosecond(),
		t.Time.Location())
}

// IsLunarLeapMonthPossible takes a lunar date that lacks leap year information
// and tries to figure out if it was possible that month was repeated.
func IsLunarLeapMonthPossible(t LunarTime) bool {
	// Take the lunar time assuming it's not a leap month
	t.IsLeap = false

	// Assuming a month can't be longer than 31 days, we jump forward by that
	// amount plus 1 to land in the next month.
	diff := 31 - t.Time.Day()
	lunarPlus := t.Add(day * time.Duration(diff+1))

	return lunarPlus.IsLeap
}

func (t LunarTime) Add(d time.Duration) LunarTime {
	return SolarToLunar(LunarToSolar(t).Add(d))
}

func (t LunarTime) Sub(u LunarTime) time.Duration {
	return LunarToSolar(t).Sub(LunarToSolar(u))
}

func (t LunarTime) Equal(u LunarTime) bool {
	return t.Time.Equal(u.Time) && t.IsLeap == u.IsLeap
}

func (t LunarTime) Before(u LunarTime) bool {
	return LunarToSolar(t).Before(LunarToSolar(u))
}

func (t LunarTime) After(u LunarTime) bool {
	return LunarToSolar(t).After(LunarToSolar(u))
}
