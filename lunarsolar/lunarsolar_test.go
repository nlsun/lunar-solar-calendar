package lunarsolar

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLunarToSolar(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		lunar    LunarTime
		expected time.Time
	}{
		{
			scenario: "not leap year",
			lunar: LunarTime{
				time:   time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC),
				isLeap: false,
			},
			expected: time.Date(2019, 4, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			// It appears that calling a non-leap month a leap month will
			// give an incorrect result.
			scenario: "not leap year, but calling it such",
			lunar: LunarTime{
				time:   time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC),
				isLeap: true,
			},
			expected: time.Date(2019, 2, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			scenario: "leap month",
			lunar: LunarTime{
				time:   time.Date(1998, 5, 2, 0, 0, 0, 0, time.UTC),
				isLeap: true,
			},
			expected: time.Date(1998, 6, 25, 0, 0, 0, 0, time.UTC),
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			solar := LunarToSolar(tc.lunar)
			assert.Equal(t, tc.expected, solar, fmt.Sprintf("%v\n%v", tc.expected, solar))
		})
	}
}

func TestSolarToLunar(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		solar    time.Time
		expected LunarTime
	}{
		{
			scenario: "not leap year",
			solar:    time.Date(2019, 4, 5, 0, 0, 0, 0, time.UTC),
			expected: LunarTime{
				time:   time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC),
				isLeap: false,
			},
		},
		{
			scenario: "is leap year, not leap month",
			solar:    time.Date(2020, 1, 26, 0, 0, 0, 0, time.UTC),
			expected: LunarTime{
				time:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				isLeap: false,
			},
		},
		{
			scenario: "is leap year, is leap month, but not the duplicate",
			solar:    time.Date(2020, 4, 23, 0, 0, 0, 0, time.UTC),
			expected: LunarTime{
				time:   time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
				isLeap: false,
			},
		},
		{
			scenario: "is leap year, is leap month",
			solar:    time.Date(2020, 5, 23, 0, 0, 0, 0, time.UTC),
			expected: LunarTime{
				time:   time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
				isLeap: true,
			},
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			lunar := SolarToLunar(tc.solar)
			assert.Equal(t, tc.expected.time, lunar.time, fmt.Sprintf("%v\n%v", tc.expected.time, lunar.time))
			assert.Equal(t, tc.expected.isLeap, lunar.isLeap)
		})
	}
}

func TestIsLunarLeapMonthPossible(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		lunar    LunarTime
		expected bool
	}{
		{
			scenario: "not leap year",
			lunar:    LunarTime{time: time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC)},
			expected: false,
		},
		{
			scenario: "not leap, month before",
			lunar:    LunarTime{time: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)},
			expected: false,
		},
		{
			scenario: "possible leap month",
			lunar:    LunarTime{time: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)},
			expected: true,
		},
		{
			scenario: "not leap, month after",
			lunar:    LunarTime{time: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC)},
			expected: false,
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			res := IsLunarLeapMonthPossible(tc.lunar)
			assert.Equal(t, tc.expected, res)
		})
	}
}
