package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"caldav-cli/internal/auth"
	"caldav-cli/internal/client"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Store CalDAV server credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Server URL: ")
		server, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading server URL: %w", err)
		}
		server = strings.TrimSpace(server)

		fmt.Print("Username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading username: %w", err)
		}
		username = strings.TrimSpace(username)

		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("reading password: %w", err)
		}
		fmt.Println()
		password := string(passwordBytes)

		creds := auth.Credentials{
			Server:   server,
			Username: username,
			Password: password,
		}

		fmt.Print("Validating credentials... ")
		c, err := client.New(creds, Verbose)
		if err != nil {
			return fmt.Errorf("connecting: %w", err)
		}
		_, err = c.ListCalendars(context.Background())
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
		fmt.Println("OK")

		if err := auth.Store(creds); err != nil {
			return fmt.Errorf("storing credentials: %w", err)
		}
		fmt.Println("Credentials saved.")
		return nil
	},
}
