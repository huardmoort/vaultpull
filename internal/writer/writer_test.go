package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteEnvFile_Basic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := WriteEnvFile(path, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	content := string(data)
	for k, v := range secrets {
		expected := k + "=" + v
		if !strings.Contains(content, expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, content)
		}
	}
}

func TestWriteEnvFile_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", ".env")

	if err := WriteEnvFile(path, map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestWriteEnvFile_EscapesSpaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := WriteEnvFile(path, map[string]string{"MSG": "hello world"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", data)
	}
}

func TestWriteEnvFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := WriteEnvFile(path, map[string]string{"X": "1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected permissions 0600, got %v", info.Mode().Perm())
	}
}
