package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"caldav-cli/internal/auth"
	"caldav-cli/internal/client"
	"caldav-cli/internal/format"

	"github.com/emersion/go-webdav/caldav"
	"github.com/spf13/cobra"
)

var (
	eventsCalendar string
	eventsFrom     string
	eventsTo       string
	eventsSearch   string
	eventsJSON     bool
	eventsRaw      bool
)

func init() {
	eventsCmd.Flags().StringVar(&eventsCalendar, "calendar", "", "Calendar path or name")
	eventsCmd.Flags().StringVar(&eventsFrom, "from", "", "Start date (YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)")
	eventsCmd.Flags().StringVar(&eventsTo, "to", "", "End date (YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)")
	eventsCmd.Flags().StringVar(&eventsSearch, "search", "", "Filter events by summary or description (case-insensitive)")
	eventsCmd.Flags().BoolVar(&eventsJSON, "json", false, "Output as JSON")
	eventsCmd.Flags().BoolVar(&eventsRaw, "raw", false, "Output raw iCal data")
	rootCmd.AddCommand(eventsCmd)
}

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List events in a date range",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.Load()
		if err != nil {
			return err
		}
		c, err := client.New(creds, Verbose)
		if err != nil {
			return err
		}
		ctx := context.Background()

		calPath, err := resolveCalendar(ctx, c, eventsCalendar)
		if err != nil {
			return err
		}

		from, to, err := parseDateRange(eventsFrom, eventsTo)
		if err != nil {
			return err
		}

		objects, err := c.ListEvents(ctx, calPath, from, to)
		if err != nil {
			return err
		}

		if eventsSearch != "" {
			objects = filterObjects(objects, eventsSearch)
		}

		return printOutput(objects, eventsJSON, eventsRaw)
	},
}

func parseDate(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format %q (use YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)", s)
}

func parseDateRange(fromStr, toStr string) (time.Time, time.Time, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var from, to time.Time
	var err error

	if fromStr == "" {
		from = today
	} else {
		from, err = parseDate(fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if toStr == "" {
		to = from.AddDate(0, 0, 30)
	} else {
		to, err = parseDate(toStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	return from, to, nil
}

func resolveCalendar(ctx context.Context, c *client.Client, flag string) (string, error) {
	calendars, err := c.ListCalendars(ctx)
	if err != nil {
		return "", fmt.Errorf("listing calendars: %w", err)
	}
	if len(calendars) == 0 {
		return "", fmt.Errorf("no calendars found")
	}

	if flag == "" {
		if len(calendars) == 1 {
			return calendars[0].Path, nil
		}
		names := make([]string, len(calendars))
		for i, cal := range calendars {
			names[i] = fmt.Sprintf("  %s (%s)", cal.Name, cal.Path)
		}
		return "", fmt.Errorf("multiple calendars found, use --calendar to specify one:\n%s",
			joinLines(names))
	}

	for _, cal := range calendars {
		if cal.Path == flag || cal.Name == flag {
			return cal.Path, nil
		}
	}
	return "", fmt.Errorf("calendar %q not found", flag)
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func filterObjects(objects []caldav.CalendarObject, search string) []caldav.CalendarObject {
	term := strings.ToLower(search)
	var filtered []caldav.CalendarObject
	for _, obj := range objects {
		if obj.Data == nil {
			continue
		}
		for _, child := range obj.Data.Children {
			if child.Name != "VEVENT" {
				continue
			}
			if prop := child.Props.Get("SUMMARY"); prop != nil {
				if strings.Contains(strings.ToLower(prop.Value), term) {
					filtered = append(filtered, obj)
					break
				}
			}
			if prop := child.Props.Get("DESCRIPTION"); prop != nil {
				if strings.Contains(strings.ToLower(prop.Value), term) {
					filtered = append(filtered, obj)
					break
				}
			}
		}
	}
	return filtered
}

func printOutput(objects []caldav.CalendarObject, jsonFlag, rawFlag bool) error {
	if rawFlag {
		return format.PrintRawICal(objects)
	}
	events := format.ExtractEvents(objects)
	if jsonFlag {
		return format.PrintJSON(events)
	}
	format.PrintTable(events)
	return nil
}
