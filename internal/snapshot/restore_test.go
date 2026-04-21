package snapshot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRestore_WritesSecrets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	entry := Entry{
		Path:      path,
		CapturedAt: time.Now(),
		Secrets:   map[string]string{"FOO": "bar", "BAZ": "qux"},
	}

	res := Restore(entry)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Written != 2 {
		t.Fatalf("expected 2 written, got %d", res.Written)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", content)
	}
	if !strings.Contains(content, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", content)
	}
}

func TestRestore_CreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", ".env")

	entry := Entry{
		Path:      path,
		CapturedAt: time.Now(),
		Secrets:   map[string]string{"KEY": "val"},
	}

	res := Restore(entry)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestRestore_EmptyPath(t *testing.T) {
	entry := Entry{Path: "", Secrets: map[string]string{"X": "y"}}
	res := Restore(entry)
	if res.Err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestRestoreAll_MultipleEntries(t *testing.T) {
	dir := t.TempDir()

	entries := []Entry{
		{Path: filepath.Join(dir, "a.env"), CapturedAt: time.Now(), Secrets: map[string]string{"A": "1"}},
		{Path: filepath.Join(dir, "b.env"), CapturedAt: time.Now(), Secrets: map[string]string{"B": "2"}},
	}

	results := RestoreAll(entries)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %s: %v", r.Path, r.Err)
		}
	}
}
