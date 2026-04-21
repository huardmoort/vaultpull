package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
)

// RestoreResult holds the outcome of a single file restore operation.
type RestoreResult struct {
	Path    string
	Written int
	Err     error
}

// Restore writes the secrets from a snapshot entry back to the target env file.
// It overwrites the file at entry.Path with the key=value pairs stored in the snapshot.
// The parent directory is created if it does not exist.
func Restore(entry Entry) RestoreResult {
	if entry.Path == "" {
		return RestoreResult{Path: entry.Path, Err: fmt.Errorf("snapshot entry has empty path")}
	}

	if err := os.MkdirAll(filepath.Dir(entry.Path), 0o755); err != nil {
		return RestoreResult{Path: entry.Path, Err: fmt.Errorf("create parent dir: %w", err)}
	}

	f, err := os.OpenFile(entry.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return RestoreResult{Path: entry.Path, Err: fmt.Errorf("open file: %w", err)}
	}
	defer f.Close()

	written := 0
	for k, v := range entry.Secrets {
		line := fmt.Sprintf("%s=%s\n", k, v)
		if _, err := fmt.Fprint(f, line); err != nil {
			return RestoreResult{Path: entry.Path, Written: written, Err: fmt.Errorf("write key %q: %w", k, err)}
		}
		written++
	}

	return RestoreResult{Path: entry.Path, Written: written}
}

// RestoreAll restores multiple snapshot entries and returns one result per entry.
func RestoreAll(entries []Entry) []RestoreResult {
	results := make([]RestoreResult, 0, len(entries))
	for _, e := range entries {
		results = append(results, Restore(e))
	}
	return results
}
