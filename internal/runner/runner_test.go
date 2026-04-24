package runner

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/lock"
)

type fakeVault struct {
	data map[string]map[string]string
	err  error
}

func (f *fakeVault) ReadSecret(path string) (map[string]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	if v, ok := f.data[path]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found: %s", path)
}

func newFakeVault(data map[string]map[string]string) *fakeVault {
	return &fakeVault{data: data}
}

func TestRun_Success(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	cfg := &config.Config{
		VaultAddr: "http://127.0.0.1:8200",
		Mappings: []config.Mapping{
			{
				EnvFile: envFile,
				Secrets: []config.SecretMapping{
					{VaultPath: "secret/app", EnvKey: "DB_URL", VaultKey: "url"},
				},
			},
		},
	}

	client := newFakeVault(map[string]map[string]string{
		"secret/app": {"url": "postgres://localhost/db"},
	})

	results, err := runWithClient(cfg, client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Errorf("expected no error in result, got: %v", results[0].Err)
	}
	if results[0].Keys != 1 {
		t.Errorf("expected 1 key written, got %d", results[0].Keys)
	}
}

func TestRun_FetchError(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	cfg := &config.Config{
		Mappings: []config.Mapping{
			{
				EnvFile: envFile,
				Secrets: []config.SecretMapping{
					{VaultPath: "secret/missing", EnvKey: "X", VaultKey: "x"},
				},
			},
		},
	}

	client := &fakeVault{err: fmt.Errorf("vault unavailable")}
	results, err := runWithClient(cfg, client)
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if results[0].Err == nil {
		t.Error("expected result error for failed fetch")
	}
}

func TestRun_LockedFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	// Pre-acquire the lock to simulate a concurrent process.
	lf, err := lock.Acquire(envFile)
	if err != nil {
		t.Fatalf("setup lock failed: %v", err)
	}
	defer lf.Release()

	cfg := &config.Config{
		Mappings: []config.Mapping{
			{EnvFile: envFile, Secrets: []config.SecretMapping{}},
		},
	}

	client := newFakeVault(nil)
	results, err := runWithClient(cfg, client)
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if results[0].Err == nil {
		t.Error("expected lock contention error in result")
	}
}
