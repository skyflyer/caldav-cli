package format

import (
	"fmt"
	"os"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

func PrintRawICal(objects []caldav.CalendarObject) error {
	enc := ical.NewEncoder(os.Stdout)
	for _, obj := range objects {
		if obj.Data == nil {
			continue
		}
		if err := enc.Encode(obj.Data); err != nil {
			return fmt.Errorf("encoding ical: %w", err)
		}
	}
	return nil
}
