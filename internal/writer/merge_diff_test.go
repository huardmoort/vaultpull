package writer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeWithDiff_NewKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}
	changes, err := MergeWithDiff(path, incoming, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	for _, c := range changes {
		if c.Action != "add" {
			t.Errorf("expected add, got %s for key %s", c.Action, c.Key)
		}
	}
}

func TestMergeWithDiff_UpdateKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// seed existing file
	if err := WriteEnvFile(path, map[string]string{"FOO": "old"}); err != nil {
		t.Fatal(err)
	}

	changes, err := MergeWithDiff(path, map[string]string{"FOO": "new"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 1 || changes[0].Action != "update" {
		t.Fatalf("expected 1 update, got %+v", changes)
	}
}

func TestMergeWithDiff_NoChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := WriteEnvFile(path, map[string]string{"FOO": "bar"}); err != nil {
		t.Fatal(err)
	}

	changes, err := MergeWithDiff(path, map[string]string{"FOO": "bar"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %+v", changes)
	}
}

func TestMergeWithDiff_VerboseNoError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_, err := MergeWithDiff(path, map[string]string{"X": "1"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}
