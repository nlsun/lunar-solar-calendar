package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nlsun/lunar-solar-calendar/lunarsolar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	day = 24 * time.Hour
)

func TestGenerateLunarBirthdayCalendar(t *testing.T) {
	for _, tc := range []struct {
		scenario      string
		lunarBirth    lunarsolar.LunarTime
		lastYear      int
		title         string
		description   string
		notifications []notification
		expected      string
	}{
		{
			// Tests multiple notifications
			// Tests multiple years
			scenario: "google calendar",
			lunarBirth: lunarsolar.NewLunarTime(
				time.Date(2020, 11, 6, 0, 0, 0, 0, time.UTC),
				false,
			),
			lastYear:    2022,
			title:       "test-title",
			description: "test-description",
			notifications: []notification{
				{Duration: jsonDuration{9 * time.Hour}, Forward: true},
				{Duration: jsonDuration{15 * time.Hour}},
				{Duration: jsonDuration{6*day + 15*time.Hour}},
				{Duration: jsonDuration{13*day + 15*time.Hour}},
			},
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			cal := generateLunarBirthdayCalendar(tc.lunarBirth, tc.lastYear, tc.title, tc.description, tc.notifications)

			err := ioutil.WriteFile("calendar_test_output.ics", []byte(cal.Serialize()), 0644)
			require.NoError(t, err)

			b, err := ioutil.ReadFile(tc.scenario + ".ics")
			require.NoError(t, err)

			assert.Equal(t, string(b), cal.Serialize())
		})
	}
}

func TestFormatDuration(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		duration notification
		expected string
	}{
		{
			scenario: "backward, 1 day, 1 hour, 1 minute, 1 second",
			duration: notification{
				Duration: jsonDuration{24*time.Hour + time.Hour + time.Minute + time.Second},
			},
			expected: "-P1DT1H1M1S",
		},
		{
			scenario: "forward, 1 day, 1 hour, 1 minute, 1 second",
			duration: notification{
				Duration: jsonDuration{24*time.Hour + time.Hour + time.Minute + time.Second},
				Forward:  true,
			},
			expected: "P1DT1H1M1S",
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			assert.Equal(t, tc.expected, formatDuration(tc.duration))
		})
	}
}

func TestJsonDurationUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		scenario string
		duration string
		expected jsonDuration
	}{
		{
			scenario: "string",
			duration: `"1h"`,
			expected: jsonDuration{time.Duration(time.Hour)},
		},
		{
			scenario: "float64",
			duration: `3600000000000`,
			expected: jsonDuration{time.Duration(time.Hour)},
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			var v jsonDuration
			err := json.Unmarshal([]byte(tc.duration), &v)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, v)
		})
	}
}
