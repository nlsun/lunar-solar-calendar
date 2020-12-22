package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
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
			scenario: "born normal year, target normal month",
			lunarBirth: lunarsolar.NewLunarTime(
				time.Date(1958, 11, 6, 0, 0, 0, 0, time.UTC),
				false,
			),
			targetSolarYear: 2020,
			expected:        time.Date(2020, 12, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			scenario: "born leap month, target leap month",
			lunarBirth: lunarsolar.NewLunarTime(
				time.Date(1998, 5, 2, 0, 0, 0, 0, time.UTC),
				true,
			),
			targetSolarYear: 2009,
			expected:        time.Date(2009, 6, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			scenario: "born leap month, target normal month",
			lunarBirth: lunarsolar.NewLunarTime(
				time.Date(1998, 5, 2, 0, 0, 0, 0, time.UTC),
				true,
			),
			targetSolarYear: 2010,
			expected:        time.Date(2010, 6, 13, 0, 0, 0, 0, time.UTC),
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			solarBirthday, err := lunarBirthdayForYear(tc.lunarBirth, tc.targetSolarYear)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, solarBirthday, fmt.Sprintf("%v\n%v", tc.expected, solarBirthday))
		})
	}
}

func TestLunarBirthdayForYearHTTP(t *testing.T) {
	s := httptest.NewServer(mkHandler(""))
	defer s.Close()

	for _, tc := range []struct {
		scenario string
		request  map[string]interface{}
		expected time.Time
	}{
		{
			scenario: "not leap year",
			request: map[string]interface{}{
				"lunar_birth_date": time.Date(1958, 11, 6, 0, 0, 0, 0, time.UTC),
				"is_leap_month":    false,
				"year":             2020,
			},
			expected: time.Date(2020, 12, 20, 0, 0, 0, 0, time.UTC),
		},
	} {
		t.Run(tc.scenario, func(t *testing.T) {
			b, err := json.Marshal(tc.request)
			require.NoError(t, err)

			reqURL := s.URL + "/api/v1/lunar-birthday-for-year/"
			resp, err := s.Client().Post(reqURL, "application/json", bytes.NewReader(b))
			require.NoError(t, err)
			defer resp.Body.Close()

			b, err = ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			var respBody lunarBirthdayForYearResponse
			err = json.Unmarshal(b, &respBody)
			require.NoError(t, err)

			date := time.Date(respBody.Year, time.Month(respBody.Month), respBody.Day, 0.0, 0, 0, 0, time.UTC)
			assert.Equal(t, tc.expected, date, fmt.Sprintf("%v\n%v", tc.expected, date))
		})
	}
}
