package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `
vault_addr: http://127.0.0.1:8200
vault_token: root
mappings:
  - vault_path: secret/data/app
    env_file: .env
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected vault_addr: %s", cfg.VaultAddr)
	}
	if len(cfg.Mappings) != 1 {
		t.Errorf("expected 1 mapping, got %d", len(cfg.Mappings))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MissingVaultAddr(t *testing.T) {
	path := writeTemp(t, `
vault_token: root
mappings:
  - vault_path: secret/data/app
    env_file: .env
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLoad_EmptyMappings(t *testing.T) {
	path := writeTemp(t, `
vault_addr: http://127.0.0.1:8200
vault_token: root
mappings: []
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty mappings")
	}
}

func TestLoad_MissingVaultToken(t *testing.T) {
	path := writeTemp(t, `
vault_addr: http://127.0.0.1:8200
mappings:
  - vault_path: secret/data/app
    env_file: .env
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing vault_token")
	}
}
