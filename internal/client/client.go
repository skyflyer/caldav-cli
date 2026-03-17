package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"caldav-cli/internal/auth"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

type Client struct {
	dav *caldav.Client
}

func New(creds auth.Credentials, verbose bool) (*Client, error) {
	var httpClient webdav.HTTPClient
	httpClient = webdav.HTTPClientWithBasicAuth(nil, creds.Username, creds.Password)
	if verbose {
		httpClient = &loggingClient{
			inner:  httpClient,
			logger: log.New(os.Stderr, "[HTTP] ", log.LstdFlags),
		}
	}
	davClient, err := caldav.NewClient(httpClient, creds.Server)
	if err != nil {
		return nil, fmt.Errorf("creating caldav client: %w", err)
	}
	return &Client{dav: davClient}, nil
}

type loggingClient struct {
	inner  webdav.HTTPClient
	logger *log.Logger
}

func (c *loggingClient) Do(req *http.Request) (*http.Response, error) {
	c.logger.Printf("--> %s %s", req.Method, req.URL)

	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("reading request body for logging: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
		if len(body) > 0 {
			c.logger.Printf("--> Body:\n%s", body)
		}
	}

	resp, err := c.inner.Do(req)
	if err != nil {
		c.logger.Printf("<-- ERROR: %v", err)
		return nil, err
	}

	c.logger.Printf("<-- %s", resp.Status)

	if resp.Body != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body for logging: %w", err)
		}
		resp.Body = io.NopCloser(bytes.NewReader(body))
		if len(body) > 0 {
			c.logger.Printf("<-- Body:\n%s", body)
		}
	}

	return resp, nil
}

func (c *Client) ListCalendars(ctx context.Context) ([]caldav.Calendar, error) {
	principal, err := c.dav.FindCurrentUserPrincipal(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding user principal: %w", err)
	}
	homeSet, err := c.dav.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		return nil, fmt.Errorf("finding calendar home set: %w", err)
	}
	calendars, err := c.dav.FindCalendars(ctx, homeSet)
	if err != nil {
		return nil, fmt.Errorf("finding calendars: %w", err)
	}
	return calendars, nil
}

func (c *Client) ListEvents(ctx context.Context, calendarPath string, from, to time.Time) ([]caldav.CalendarObject, error) {
	query := &caldav.CalendarQuery{
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{
				{
					Name:  "VEVENT",
					Start: from,
					End:   to,
				},
			},
		},
	}
	objects, err := c.dav.QueryCalendar(ctx, calendarPath, query)
	if err != nil {
		return nil, fmt.Errorf("querying events: %w", err)
	}
	return objects, nil
}

func (c *Client) GetEvent(ctx context.Context, calendarPath string, uid string) ([]caldav.CalendarObject, error) {
	query := &caldav.CalendarQuery{
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{
				{
					Name: "VEVENT",
					Props: []caldav.PropFilter{
						{
							Name:      "UID",
							TextMatch: &caldav.TextMatch{Text: uid},
						},
					},
				},
			},
		},
	}
	objects, err := c.dav.QueryCalendar(ctx, calendarPath, query)
	if err != nil {
		return nil, fmt.Errorf("querying event %s: %w", uid, err)
	}
	return objects, nil
}
