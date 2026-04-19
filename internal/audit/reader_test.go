package audit_test

import (
	"os"
	"testing"

	"github.com/yourorg/vaultpull/internal/audit"
)

func TestReadAll_Empty(t *testing.T) {
	tmp, _ := os.CreateTemp(t.TempDir(), "audit-*.jsonl")
	tmp.Close()

	entries, err := audit.ReadAll(tmp.Name())
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestReadAll_RoundTrip(t *testing.T) {
	tmp, _ := os.CreateTemp(t.TempDir(), "audit-*.jsonl")
	tmp.Close()

	l, _ := audit.New(tmp.Name())
	_ = l.Log("sync", "secret/x", ".env", "ok", "all good")
	_ = l.Log("sync", "secret/y", ".env.prod", "error", "forbidden")
	l.Close()

	entries, err := audit.ReadAll(tmp.Name())
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Path != "secret/x" {
		t.Errorf("entry[0].Path = %s", entries[0].Path)
	}
	if entries[1].Status != "error" {
		t.Errorf("entry[1].Status = %s", entries[1].Status)
	}
	if entries[1].Message != "forbidden" {
		t.Errorf("entry[1].Message = %s", entries[1].Message)
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	_, err := audit.ReadAll("/tmp/does-not-exist-vaultpull-audit.jsonl")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
