package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/vaultpull/internal/vault"
)

func newFakeVault(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": payload})
	}))
}

func TestReadSecret_KVv1(t *testing.T) {
	expected := map[string]interface{}{"API_KEY": "abc123"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": expected})
	}))
	defer srv.Close()

	c, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	data, err := c.ReadSecret("secret/myapp")
	if err != nil {
		t.Fatalf("ReadSecret: %v", err)
	}

	if data["API_KEY"] != "abc123" {
		t.Errorf("expected abc123, got %v", data["API_KEY"])
	}
}

func TestReadSecret_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.ReadSecret("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}
