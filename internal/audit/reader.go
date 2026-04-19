package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// ReadAll reads all audit entries from the log file at path.
func ReadAll(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("audit: open for read: %w", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse line: %w", err)
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("audit: scan: %w", err)
	}
	return entries, nil
}
