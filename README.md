# caldav-cli

A command-line tool for browsing self-hosted CalDAV servers (Stalwart, Radicale, Baikal, etc.).

Credentials are stored in the OS keychain. Output defaults to a human-readable table, with `--json` and `--raw` flags for machine-readable and iCal formats.

## Install

```
go install caldav-cli@latest
```

Or build from source:

```
go build -o caldav-cli .
```

## Usage

### Login

Prompts for server URL, username, and password. Validates credentials by connecting to the server before storing them in the OS keychain.

```
caldav-cli login
```

### List calendars

```
caldav-cli calendars
```

### List events

```
# Events from today + 30 days (default)
caldav-cli events

# Specific date range
caldav-cli events --from 2024-01-01 --to 2024-12-31

# Filter by calendar (path or name)
caldav-cli events --calendar personal

# JSON output
caldav-cli events --json

# Raw iCal output
caldav-cli events --raw
```

Dates accept `YYYY-MM-DD` and `YYYY-MM-DDTHH:MM:SS` formats. When only `--from` is given, `--to` defaults to 30 days later.

### Fetch a single event

```
# Searches all calendars
caldav-cli event <uid>

# Specific calendar
caldav-cli event <uid> --calendar personal

# Raw iCal
caldav-cli event <uid> --raw
```

### Logout

Removes stored credentials from the OS keychain.

```
caldav-cli logout
```

## Dependencies

| Library | Purpose |
|---------|---------|
| [cobra](https://github.com/spf13/cobra) | CLI framework |
| [go-webdav](https://github.com/emersion/go-webdav) | CalDAV client |
| [go-ical](https://github.com/emersion/go-ical) | iCal encoding |
| [go-keyring](https://github.com/zalando/go-keyring) | OS keychain storage |
| [table](https://github.com/rodaine/table) | Table output |
| [x/term](https://pkg.go.dev/golang.org/x/term) | Hidden password input |
