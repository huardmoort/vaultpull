package audit

import (
	"strings"
	"testing"
	"time"
)

func makeEntry(ts time.Time, added, updated, deleted int, errMsg string) Entry {
	return Entry{
		Timestamp: ts,
		Added:     added,
		Updated:   updated,
		Deleted:   deleted,
		Error:     errMsg,
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.TotalRuns != 0 {
		t.Errorf("expected 0 total runs, got %d", s.TotalRuns)
	}
	if s.FirstRun != nil || s.LastRun != nil {
		t.Error("expected nil FirstRun and LastRun for empty entries")
	}
}

func TestSummarize_Counts(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)
	entries := []Entry{
		makeEntry(earlier, 3, 1, 0, ""),
		makeEntry(now, 0, 2, 1, ""),
		makeEntry(now.Add(-30*time.Minute), 1, 0, 0, "vault unreachable"),
	}
	s := Summarize(entries)
	if s.TotalRuns != 3 {
		t.Errorf("expected 3 total runs, got %d", s.TotalRuns)
	}
	if s.SuccessRuns != 2 {
		t.Errorf("expected 2 success runs, got %d", s.SuccessRuns)
	}
	if s.FailedRuns != 1 {
		t.Errorf("expected 1 failed run, got %d", s.FailedRuns)
	}
	if s.SecretsAdded != 4 {
		t.Errorf("expected 4 added, got %d", s.SecretsAdded)
	}
	if s.SecretsUpdated != 3 {
		t.Errorf("expected 3 updated, got %d", s.SecretsUpdated)
	}
	if s.SecretsDeleted != 1 {
		t.Errorf("expected 1 deleted, got %d", s.SecretsDeleted)
	}
	if s.FirstRun == nil || !s.FirstRun.Equal(earlier) {
		t.Error("expected FirstRun to be the earliest timestamp")
	}
	if s.LastRun == nil || !s.LastRun.Equal(now) {
		t.Error("expected LastRun to be the latest timestamp")
	}
}

func TestFormat_Empty(t *testing.T) {
	s := Summarize(nil)
	out := s.Format()
	if !strings.Contains(out, "No audit entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormat_NonEmpty(t *testing.T) {
	entries := []Entry{
		makeEntry(time.Now(), 2, 1, 0, ""),
	}
	s := Summarize(entries)
	out := s.Format()
	if !strings.Contains(out, "Total runs:") {
		t.Errorf("expected Total runs in output, got: %s", out)
	}
	if !strings.Contains(out, "Secrets added:") {
		t.Errorf("expected Secrets added in output, got: %s", out)
	}
}
