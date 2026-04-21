package snapshot_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func makeEntry(path string, secrets map[string]string) snapshot.Entry {
	return snapshot.Entry{Path: path, Secrets: secrets}
}

func TestDiffEntries_Added(t *testing.T) {
	old := makeEntry(".env", map[string]string{})
	new := makeEntry(".env", map[string]string{"NEW_KEY": "val"})

	changes := snapshot.DiffEntries(old, new)
	if len(changes) != 1 || changes[0].Kind != snapshot.Added {
		t.Errorf("expected 1 added change, got %+v", changes)
	}
}

func TestDiffEntries_Removed(t *testing.T) {
	old := makeEntry(".env", map[string]string{"OLD_KEY": "val"})
	new := makeEntry(".env", map[string]string{})

	changes := snapshot.DiffEntries(old, new)
	if len(changes) != 1 || changes[0].Kind != snapshot.Removed {
		t.Errorf("expected 1 removed change, got %+v", changes)
	}
}

func TestDiffEntries_Changed(t *testing.T) {
	old := makeEntry(".env", map[string]string{"KEY": "old"})
	new := makeEntry(".env", map[string]string{"KEY": "new"})

	changes := snapshot.DiffEntries(old, new)
	if len(changes) != 1 || changes[0].Kind != snapshot.Changed {
		t.Errorf("expected 1 changed entry, got %+v", changes)
	}
	if changes[0].OldVal != "old" || changes[0].NewVal != "new" {
		t.Errorf("unexpected values: %+v", changes[0])
	}
}

func TestDiffEntries_NoChange(t *testing.T) {
	old := makeEntry(".env", map[string]string{"KEY": "same"})
	new := makeEntry(".env", map[string]string{"KEY": "same"})

	changes := snapshot.DiffEntries(old, new)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %+v", changes)
	}
}

func TestDiffEntries_SortedByKey(t *testing.T) {
	old := makeEntry(".env", map[string]string{})
	new := makeEntry(".env", map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"})

	changes := snapshot.DiffEntries(old, new)
	if changes[0].Key != "A_KEY" || changes[1].Key != "M_KEY" || changes[2].Key != "Z_KEY" {
		t.Errorf("changes not sorted: %+v", changes)
	}
}

func TestFormatDiff_Empty(t *testing.T) {
	out := snapshot.FormatDiff(nil)
	if out != "no changes" {
		t.Errorf("expected 'no changes', got %q", out)
	}
}

func TestFormatDiff_NonEmpty(t *testing.T) {
	old := makeEntry(".env", map[string]string{"OLD": "v", "SAME": "x"})
	new := makeEntry(".env", map[string]string{"NEW": "v", "SAME": "y"})

	changes := snapshot.DiffEntries(old, new)
	out := snapshot.FormatDiff(changes)

	if !strings.Contains(out, "+ NEW") {
		t.Errorf("expected '+ NEW' in output: %s", out)
	}
	if !strings.Contains(out, "- OLD") {
		t.Errorf("expected '- OLD' in output: %s", out)
	}
	if !strings.Contains(out, "~ SAME") {
		t.Errorf("expected '~ SAME' in output: %s", out)
	}
}
