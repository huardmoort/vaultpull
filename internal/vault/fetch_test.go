package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/vault"
)

func TestFetchAll(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"DB_PASS": "secret"},
		})
	}))
	defer srv.Close()

	c, err := vault.NewClient(srv.URL, "tok")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	mappings := []config.Mapping{
		{VaultPath: "secret/db", EnvFile: ".env"},
	}

	result, err := vault.FetchAll(c, mappings)
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}

	sm, ok := result[".env"]
	if !ok {
		t.Fatal("expected .env key in result")
	}
	if sm["DB_PASS"] != "secret" {
		t.Errorf("expected 'secret', got %q", sm["DB_PASS"])
	}
}

func TestFetchAll_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c, err := vault.NewClient(srv.URL, "bad-tok")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = vault.FetchAll(c, []config.Mapping{{VaultPath: "secret/x", EnvFile: ".env"}})
	if err == nil {
		t.Fatal("expected error")
	}
}
