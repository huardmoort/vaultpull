package runner_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vaultpull/internal/config"
	"github.com/vaultpull/internal/runner"
)

func newFakeVault(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}`))
	}))
}

func TestRun_Success(t *testing.T) {
	srv := newFakeVault(t)
	defer srv.Close()

	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	cfg := &config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		Mappings: []config.Mapping{
			{
				EnvFile: envFile,
				Secrets: []config.SecretRef{
					{Path: "secret/myapp", Key: "API_KEY", EnvVar: "API_KEY"},
				},
			},
		},
	}

	results := runner.Run(cfg)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Written == 0 {
		t.Error("expected at least one secret written")
	}

	data, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("env file not created: %v", err)
	}
	if !strings.Contains(string(data), "API_KEY") {
		t.Error("env file missing API_KEY")
	}
}

func TestRun_BadVaultAddr(t *testing.T) {
	cfg := &config.Config{
		VaultAddr:  "http://127.0.0.1:0",
		VaultToken: "x",
		Mappings:   []config.Mapping{{EnvFile: "/tmp/x.env"}},
	}
	results := runner.Run(cfg)
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
}
