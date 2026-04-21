package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func tempSnapshotPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "snapshots", "state.json")
}

func TestSave_And_LoadAll(t *testing.T) {
	path := tempSnapshotPath(t)

	entry := snapshot.Entry{
		Path:    ".env",
		Secrets: map[string]string{"FOO": "bar", "BAZ": "qux"},
	}

	if err := snapshot.Save(path, entry); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := snapshot.LoadAll(path)
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", entries[0].Secrets["FOO"])
	}
}

func TestSave_ReplacesExistingPath(t *testing.T) {
	path := tempSnapshotPath(t)

	first := snapshot.Entry{Path: ".env", Secrets: map[string]string{"KEY": "old"}}
	second := snapshot.Entry{Path: ".env", Secrets: map[string]string{"KEY": "new"}}

	_ = snapshot.Save(path, first)
	_ = snapshot.Save(path, second)

	entries, err := snapshot.LoadAll(path)
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after replace, got %d", len(entries))
	}
	if entries[0].Secrets["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %s", entries[0].Secrets["KEY"])
	}
}

func TestSave_MultipleDistinctPaths(t *testing.T) {
	path := tempSnapshotPath(t)

	_ = snapshot.Save(path, snapshot.Entry{Path: ".env", Secrets: map[string]string{"A": "1"}})
	_ = snapshot.Save(path, snapshot.Entry{Path: ".env.prod", Secrets: map[string]string{"B": "2"}})

	entries, _ := snapshot.LoadAll(path)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestGetByPath_Found(t *testing.T) {
	path := tempSnapshotPath(t)
	_ = snapshot.Save(path, snapshot.Entry{Path: ".env", Secrets: map[string]string{"X": "y"}})

	e, ok, err := snapshot.GetByPath(path, ".env")
	if err != nil || !ok {
		t.Fatalf("expected found entry, ok=%v err=%v", ok, err)
	}
	if e.Secrets["X"] != "y" {
		t.Errorf("unexpected value: %s", e.Secrets["X"])
	}
}

func TestGetByPath_Missing(t *testing.T) {
	path := tempSnapshotPath(t)
	_, ok, err := snapshot.GetByPath(path, ".env")
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if ok {
		t.Error("expected not found")
	}
}

func TestLoadAll_MissingFile(t *testing.T) {
	_, err := snapshot.LoadAll(filepath.Join(t.TempDir(), "nonexistent.json"))
	if !os.IsNotExist(err) {
		t.Errorf("expected IsNotExist, got %v", err)
	}
}
