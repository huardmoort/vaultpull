package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry holds a saved snapshot of secrets for a single env file mapping.
type Entry struct {
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
	CreatedAt time.Time         `json:"created_at"`
}

// Save writes a snapshot entry to the given snapshot file (JSON, one entry per file).
func Save(snapshotPath string, entry Entry) error {
	entry.CreatedAt = time.Now().UTC()

	if err := os.MkdirAll(filepath.Dir(snapshotPath), 0o700); err != nil {
		return fmt.Errorf("snapshot: create dir: %w", err)
	}

	existing, err := LoadAll(snapshotPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot: read existing: %w", err)
	}

	// Replace or append entry for this path.
	updated := make([]Entry, 0, len(existing)+1)
	replaced := false
	for _, e := range existing {
		if e.Path == entry.Path {
			updated = append(updated, entry)
			replaced = true
		} else {
			updated = append(updated, e)
		}
	}
	if !replaced {
		updated = append(updated, entry)
	}

	f, err := os.OpenFile(snapshotPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(updated)
}

// LoadAll reads all snapshot entries from snapshotPath.
func LoadAll(snapshotPath string) ([]Entry, error) {
	f, err := os.Open(snapshotPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return entries, nil
}

// GetByPath returns the snapshot entry for a specific env file path, if present.
func GetByPath(snapshotPath, envPath string) (Entry, bool, error) {
	entries, err := LoadAll(snapshotPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Entry{}, false, nil
		}
		return Entry{}, false, err
	}
	for _, e := range entries {
		if e.Path == envPath {
			return e, true, nil
		}
	}
	return Entry{}, false, nil
}
