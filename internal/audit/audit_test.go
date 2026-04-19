package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/yourorg/vaultpull/internal/audit"
)

func TestLog_WritesEntry(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "audit-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	l, err := audit.New(tmp.Name())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer l.Close()

	if err := l.Log("sync", "secret/app", ".env", "ok", ""); err != nil {
		t.Fatalf("Log: %v", err)
	}
	l.Close()

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected one line in audit log")
	}
	var entry audit.Entry
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if entry.Operation != "sync" {
		t.Errorf("expected operation=sync, got %s", entry.Operation)
	}
	if entry.Status != "ok" {
		t.Errorf("expected status=ok, got %s", entry.Status)
	}
	if entry.Path != "secret/app" {
		t.Errorf("expected path=secret/app, got %s", entry.Path)
	}
}

func TestLog_MultipleEntries(t *testing.T) {
	tmp, _ := os.CreateTemp(t.TempDir(), "audit-*.jsonl")
	tmp.Close()

	l, _ := audit.New(tmp.Name())
	_ = l.Log("sync", "secret/a", ".env.a", "ok", "")
	_ = l.Log("sync", "secret/b", ".env.b", "error", "not found")
	l.Close()

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 entries, got %d", count)
	}
}

func TestNew_BadPath(t *testing.T) {
	_, err := audit.New("/nonexistent/dir/audit.log")
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}
