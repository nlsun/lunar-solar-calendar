package main

import (
	"encoding/json"
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/nlsun/lunar-solar-calendar/lunarsolar"
)

type vAlarm struct {
	ics.VAlarm
}

type jsonDuration struct {
	time.Duration
}

// https://www.kanzaki.com/docs/ical/duration-t.html
type notification struct {
	Duration jsonDuration `json:"duration"`
	// If true, sets the notification forward in time, otherwise backwards
	Forward bool `json:"forward"`
}

// There's a really strange bug where VALARM isn't recognized by Google
// Calendar. Even if you export a Google Calendar and re-import it to a fresh
// Google Calendar it won't work.
func generateLunarBirthdayCalendar(birthDate lunarsolar.LunarTime, lastYear int, title, description string, notifications []notification) *ics.Calendar {
	cal := ics.NewCalendar()

	for d := birthDate.Add(0); d.Time().Year() <= lastYear; d = d.AddDate(1, 0, 0) {
		birthday := lunarsolar.LunarToSolar(d)

		ev := cal.AddEvent(fmt.Sprintf("%s-%v", title, birthday))
		ev.SetSummary(title)
		ev.SetDescription(description)
		ev.SetAllDayStartAt(birthday)
		for _, notif := range notifications {
			am := addVAlarm(ev)
			am.setTrigger(notif)
		}
	}

	return cal
}

// Alarm configured to send a notification
func addVAlarm(event *ics.VEvent) *vAlarm {
	alarm := &vAlarm{}
	alarm.setProperty("ACTION", "DISPLAY")
	event.Components = append(event.Components, alarm)
	return alarm
}

// How long before event to set the alarm
func (alarm *vAlarm) setTrigger(d notification, props ...ics.PropertyParameter) {
	alarm.setProperty("TRIGGER", formatDuration(d), props...)
	// DESCRIPTION is required
	alarm.setProperty("DESCRIPTION", "This is an event reminder")
}

// Only supports time-based durations that count backwards. And of that, only
// support the nearest hour rounded down.
//
// https://www.kanzaki.com/docs/ical/duration-t.html
func formatDuration(d notification) string {
	allHours := int(d.Duration.Hours())
	allMinutes := int(d.Duration.Minutes())
	allSeconds := int(d.Duration.Seconds())

	days := allHours / 24
	hours := allHours % 24
	minutes := allMinutes % 60
	seconds := allSeconds % 60

	var direction string
	if !d.Forward {
		direction = "-"
	}
	return fmt.Sprintf("%sP%dDT%dH%dM%dS", direction, days, hours, minutes, seconds)
}

// Copied from `ics.VEvent.SetProperty`
func (alarm *vAlarm) setProperty(property ics.ComponentProperty, value string, props ...ics.PropertyParameter) {
	for i := range alarm.Properties {
		if alarm.Properties[i].IANAToken == string(property) {
			alarm.Properties[i].Value = value
			alarm.Properties[i].ICalParameters = map[string][]string{}
			for _, p := range props {
				k, v := p.KeyValue()
				alarm.Properties[i].ICalParameters[k] = v
			}
			return
		}
	}
	alarm.AddProperty(property, value, props...)
}

// Copied from `ics.VEvent.AddProperty`
func (alarm *vAlarm) AddProperty(property ics.ComponentProperty, value string, props ...ics.PropertyParameter) {
	r := ics.IANAProperty{
		BaseProperty: ics.BaseProperty{
			IANAToken:      string(property),
			Value:          value,
			ICalParameters: map[string][]string{},
		},
	}
	for _, p := range props {
		k, v := p.KeyValue()
		r.ICalParameters[k] = v
	}
	alarm.Properties = append(alarm.Properties, r)
}

// Override default behavior which only accepts an int
func (d *jsonDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = jsonDuration{time.Duration(value)}
		return nil
	case string:
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = jsonDuration{time.Duration(duration)}
		return nil
	default:
		return fmt.Errorf("invalid duration: %s", string(b))
	}
}
