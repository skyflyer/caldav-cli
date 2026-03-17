package format

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(events []EventData) error {
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		return fmt.Errorf("writing json: %w", err)
	}
	fmt.Println()
	return nil
}
