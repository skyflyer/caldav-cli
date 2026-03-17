package format

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/rodaine/table"
)

type EventData struct {
	UID         string    `json:"uid"`
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	Start       string    `json:"start"`
	End         string    `json:"end"`
	Duration    string    `json:"duration"`
	Location    string    `json:"location"`
	startTime   time.Time // parsed DTSTART, used for sorting
}

func ExtractEvents(objects []caldav.CalendarObject) []EventData {
	var events []EventData
	for _, obj := range objects {
		if obj.Data == nil {
			continue
		}
		for _, child := range obj.Data.Children {
			if child.Name != "VEVENT" {
				continue
			}
			if child.Props.Get("RECURRENCE-ID") != nil {
				continue
			}
			ev := EventData{}
			ev.UID = child.Props.Get("UID").Value
			if prop := child.Props.Get("SUMMARY"); prop != nil {
				ev.Summary = prop.Value
			}
			if prop := child.Props.Get("DESCRIPTION"); prop != nil {
				ev.Description = prop.Value
			}
			if prop := child.Props.Get("LOCATION"); prop != nil {
				ev.Location = prop.Value
			}

			var startTime, endTime time.Time
			if prop := child.Props.Get("DTSTART"); prop != nil {
				ev.Start = prop.Value
				startTime, _ = parseICal(prop)
				ev.startTime = startTime
			}
			if prop := child.Props.Get("DTEND"); prop != nil {
				ev.End = prop.Value
				endTime, _ = parseICal(prop)
			}

			if !startTime.IsZero() && !endTime.IsZero() {
				ev.Duration = formatDuration(endTime.Sub(startTime))
			}
			events = append(events, ev)
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].startTime.Before(events[j].startTime)
	})
	return events
}

func parseICal(prop *ical.Prop) (time.Time, error) {
	val := prop.Value
	layouts := []string{
		"20060102T150405Z",
		"20060102T150405",
		"20060102",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, val); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time: %s", val)
}

func formatDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	hours := d.Hours()
	if hours < 24 {
		if math.Mod(hours, 1) == 0 {
			return fmt.Sprintf("%dh", int(hours))
		}
		return fmt.Sprintf("%dh%dm", int(hours), int(d.Minutes())%60)
	}
	days := int(hours / 24)
	remainingHours := int(hours) % 24
	if remainingHours == 0 {
		return fmt.Sprintf("%dd", days)
	}
	return fmt.Sprintf("%dd%dh", days, remainingHours)
}

func PrintTable(events []EventData) {
	tbl := table.New("UID", "Summary", "Start", "End", "Duration", "Location").WithWriter(os.Stdout)
	for _, e := range events {
		tbl.AddRow(e.UID, e.Summary, e.Start, e.End, e.Duration, e.Location)
	}
	tbl.Print()
}
