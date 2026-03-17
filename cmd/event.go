package cmd

import (
	"context"
	"fmt"

	"caldav-cli/internal/auth"
	"caldav-cli/internal/client"

	"github.com/emersion/go-webdav/caldav"
	"github.com/spf13/cobra"
)

var (
	eventCalendar string
	eventJSON     bool
	eventRaw      bool
)

func init() {
	eventCmd.Flags().StringVar(&eventCalendar, "calendar", "", "Calendar path or name")
	eventCmd.Flags().BoolVar(&eventJSON, "json", false, "Output as JSON")
	eventCmd.Flags().BoolVar(&eventRaw, "raw", false, "Output raw iCal data")
	rootCmd.AddCommand(eventCmd)
}

var eventCmd = &cobra.Command{
	Use:   "event <uid>",
	Short: "Fetch a single event by UID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		uid := args[0]

		creds, err := auth.Load()
		if err != nil {
			return err
		}
		c, err := client.New(creds, Verbose)
		if err != nil {
			return err
		}
		ctx := context.Background()

		var objects []caldav.CalendarObject

		if eventCalendar != "" {
			calPath, err := resolveCalendar(ctx, c, eventCalendar)
			if err != nil {
				return err
			}
			objects, err = c.GetEvent(ctx, calPath, uid)
			if err != nil {
				return err
			}
		} else {
			calendars, err := c.ListCalendars(ctx)
			if err != nil {
				return err
			}
			for _, cal := range calendars {
				found, err := c.GetEvent(ctx, cal.Path, uid)
				if err != nil {
					continue
				}
				objects = append(objects, found...)
			}
		}

		if len(objects) == 0 {
			return fmt.Errorf("event %q not found", uid)
		}

		return printOutput(objects, eventJSON, eventRaw)
	},
}
