package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const service = "caldav-cli"

type Credentials struct {
	Server   string
	Username string
	Password string
}

func Store(creds Credentials) error {
	if err := keyring.Set(service, "server", creds.Server); err != nil {
		return fmt.Errorf("storing server: %w", err)
	}
	if err := keyring.Set(service, "username", creds.Username); err != nil {
		return fmt.Errorf("storing username: %w", err)
	}
	if err := keyring.Set(service, "password", creds.Password); err != nil {
		return fmt.Errorf("storing password: %w", err)
	}
	return nil
}

func Load() (Credentials, error) {
	server, err := keyring.Get(service, "server")
	if err != nil {
		return Credentials{}, fmt.Errorf("not logged in — run 'caldav-cli login' first")
	}
	username, err := keyring.Get(service, "username")
	if err != nil {
		return Credentials{}, fmt.Errorf("not logged in — run 'caldav-cli login' first")
	}
	password, err := keyring.Get(service, "password")
	if err != nil {
		return Credentials{}, fmt.Errorf("not logged in — run 'caldav-cli login' first")
	}
	return Credentials{
		Server:   server,
		Username: username,
		Password: password,
	}, nil
}

func Clear() error {
	var firstErr error
	for _, key := range []string{"server", "username", "password"} {
		if err := keyring.Delete(service, key); err != nil && firstErr == nil {
			if err != keyring.ErrNotFound {
				firstErr = fmt.Errorf("removing %s: %w", key, err)
			}
		}
	}
	return firstErr
}
