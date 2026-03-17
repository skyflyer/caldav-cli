package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "caldav-cli",
	Short: "Browse CalDAV calendars from the command line",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Log HTTP requests and responses")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
