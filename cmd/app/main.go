package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nlsun/lunar-solar-calendar/lunarsolar"
)

// TODO: Take traditional chinese birth date, convert to birthday of this
// Gregorian year.
// TODO: Generate google calendar for a person, notifications configurable.

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	birthday, err := lunarBirthdayForYear(lunarsolar.LunarTime{
		Time:   time.Date(1958, 11, 6, 0, 0, 0, 0, time.UTC),
		IsLeap: false,
	}, time.Now().Year())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", birthday)
}

// Computes the solar birthday given a birthday as a lunar date
// and a solar year to calculate for.
func lunarBirthdayForYear(birthDate lunarsolar.LunarTime, solarYear int) (time.Time, error) {
	solarBirth := lunarsolar.LunarToSolar(birthDate)
	solarBirthYear := solarBirth.Year()
	if solarBirthYear > solarYear {
		return time.Time{}, fmt.Errorf("birth year %d can't be greater than input year %d", solarBirthYear, solarYear)
	}

	yearDiff := solarYear - solarBirthYear
	return solarBirth.AddDate(yearDiff, 0, 0), nil
}
