package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeEntries(n int) []Entry {
	out := make([]Entry, n)
	for i := range out {
		out[i] = Entry{
			Timestamp: time.Now().UTC(),
			Mapping:   "secret/app",
			EnvFile:   ".env",
			Keys:      []string{"KEY"},
			Status:    "ok",
		}
	}
	return out
}

func writeTempEntries(t *testing.T, path string, n int) {
	t.Helper()
	if err := writeEntries(path, makeEntries(n)); err != nil {
		t.Fatalf("writeTempEntries: %v", err)
	}
}

func TestRotate_BelowThreshold(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	writeTempEntries(t, path, 10)

	archived, archivePath, err := Rotate(path, RotateOptions{MaxEntries: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if archived != 0 || archivePath != "" {
		t.Errorf("expected no rotation, got archived=%d path=%q", archived, archivePath)
	}
}

func TestRotate_ArchivesOldEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	writeTempEntries(t, path, 20)

	archived, archivePath, err := Rotate(path, RotateOptions{MaxEntries: 15, ArchiveSuffix: ".bak"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if archived != 5 {
		t.Errorf("expected 5 archived, got %d", archived)
	}
	if archivePath != path+".bak" {
		t.Errorf("unexpected archive path: %q", archivePath)
	}

	kept, err := ReadAll(path)
	if err != nil {
		t.Fatalf("read kept: %v", err)
	}
	if len(kept) != 15 {
		t.Errorf("expected 15 kept entries, got %d", len(kept))
	}

	oldEntries, err := ReadAll(archivePath)
	if err != nil {
		t.Fatalf("read archive: %v", err)
	}
	if len(oldEntries) != 5 {
		t.Errorf("expected 5 archive entries, got %d", len(oldEntries))
	}
}

func TestRotate_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.log")

	archived, archivePath, err := Rotate(path, RotateOptions{})
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if archived != 0 || archivePath != "" {
		t.Errorf("expected no-op for missing file")
	}
}

func TestRotate_DefaultSuffixIsTimestamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	writeTempEntries(t, path, 10)

	_, archivePath, err := Rotate(path, RotateOptions{MaxEntries: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(archivePath); statErr != nil {
		t.Errorf("archive file not found at %q: %v", archivePath, statErr)
	}
}
