package cmd

import (
	"fmt"

	"caldav-cli/internal/auth"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.Clear(); err != nil {
			return fmt.Errorf("clearing credentials: %w", err)
		}
		fmt.Println("Credentials removed.")
		return nil
	},
}
