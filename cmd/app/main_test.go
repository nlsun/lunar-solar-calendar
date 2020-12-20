package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/nlsun/lunar-solar-calendar/lunarsolar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLunarBirthdayForYear(t *testing.T) {
	for _, tc := range []struct {
		scenario        string
		lunarBirth      lunarsolar.LunarTime
		targetSolarYear int
		expected        time.Time
	}{
		{
			scenario: "not leap year",
			lunarBirth: lunarsolar.LunarTime{
				Time:   time.Date(1958, 11, 6, 0, 0, 0, 0, time.UTC),
				IsLeap: false,
			},
			targetSolarYear: 2020,
			expected:        time.Date(2020, 12, 16, 0, 0, 0, 0, time.UTC),
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			solarBirthday, err := lunarBirthdayForYear(tc.lunarBirth, tc.targetSolarYear)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, solarBirthday, fmt.Sprintf("%v\n%v", tc.expected, solarBirthday))
		})
	}
}
