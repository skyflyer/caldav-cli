package cmd

import (
	"context"
	"fmt"
	"os"

	"caldav-cli/internal/auth"
	"caldav-cli/internal/client"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(calendarsCmd)
}

var calendarsCmd = &cobra.Command{
	Use:   "calendars",
	Short: "List available calendars",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.Load()
		if err != nil {
			return err
		}
		c, err := client.New(creds, Verbose)
		if err != nil {
			return err
		}
		calendars, err := c.ListCalendars(context.Background())
		if err != nil {
			return err
		}

		tbl := table.New("Name", "Path", "Description").WithWriter(os.Stdout)
		for _, cal := range calendars {
			tbl.AddRow(cal.Name, cal.Path, cal.Description)
		}
		tbl.Print()
		fmt.Println()
		return nil
	},
}
