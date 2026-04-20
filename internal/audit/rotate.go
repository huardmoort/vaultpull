package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions controls how audit log rotation behaves.
type RotateOptions struct {
	// MaxEntries is the maximum number of entries to keep after rotation.
	// If zero, defaults to 500.
	MaxEntries int
	// ArchiveSuffix is appended to the original path to form the archive name.
	// If empty, defaults to "." + RFC3339 timestamp.
	ArchiveSuffix string
}

// Rotate reads the audit log at path, writes entries beyond MaxEntries to an
// archive file, and rewrites the original file with only the newest entries.
// It returns the number of entries archived and the archive path (empty if
// nothing was archived).
func Rotate(path string, opts RotateOptions) (archived int, archivePath string, err error) {
	if opts.MaxEntries <= 0 {
		opts.MaxEntries = 500
	}

	entries, err := ReadAll(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, "", nil
		}
		return 0, "", fmt.Errorf("rotate: read: %w", err)
	}

	if len(entries) <= opts.MaxEntries {
		return 0, "", nil
	}

	suffix := opts.ArchiveSuffix
	if suffix == "" {
		suffix = "." + time.Now().UTC().Format("20060102T150405Z")
	}
	archivePath = path + suffix

	old := entries[:len(entries)-opts.MaxEntries]
	keep := entries[len(entries)-opts.MaxEntries:]

	if err := writeEntries(archivePath, old); err != nil {
		return 0, "", fmt.Errorf("rotate: write archive: %w", err)
	}
	if err := writeEntries(path, keep); err != nil {
		return 0, "", fmt.Errorf("rotate: rewrite log: %w", err)
	}

	return len(old), archivePath, nil
}

func writeEntries(path string, entries []Entry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}
