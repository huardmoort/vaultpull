package writer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeEnvFile_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}

	if err := MergeEnvFile(path, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	content := string(data)
	if !contains(content, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", content)
	}
	if !contains(content, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", content)
	}
}

func TestMergeEnvFile_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	initial := "FOO=old\nKEEP=me\n"
	if err := os.WriteFile(path, []byte(initial), 0600); err != nil {
		t.Fatalf("setup error: %v", err)
	}

	secrets := map[string]string{
		"FOO": "new",
	}

	if err := MergeEnvFile(path, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !contains(content, "FOO=new") {
		t.Errorf("expected FOO=new, got:\n%s", content)
	}
	if !contains(content, "KEEP=me") {
		t.Errorf("expected KEEP=me to be preserved, got:\n%s", content)
	}
	if contains(content, "FOO=old") {
		t.Errorf("expected FOO=old to be overwritten, got:\n%s", content)
	}
}

func TestMergeEnvFile_PreservesComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	initial := "# my comment\nFOO=bar\n"
	os.WriteFile(path, []byte(initial), 0600)

	if err := MergeEnvFile(path, map[string]string{"NEW": "val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !contains(string(data), "# my comment") {
		t.Errorf("expected comment to be preserved")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
