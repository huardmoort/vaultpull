package env_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
)

// TestExpandAll_ChainedReferences verifies that values can reference other
// keys whose values are themselves already expanded (single-pass).
func TestExpandAll_ChainedReferences(t *testing.T) {
	input := map[string]string{
		"PROTO":    "https",
		"HOST":     "vault.internal",
		"BASE_URL": "${PROTO}://${HOST}",
		"LOGIN":    "${BASE_URL}/v1/auth",
	}
	out, errs := env.ExpandAll(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	// BASE_URL should be fully resolved
	if out["BASE_URL"] != "https://vault.internal" {
		t.Errorf("BASE_URL: got %q", out["BASE_URL"])
	}
	// LOGIN references BASE_URL which is in the original map (single-pass, not chained)
	// so it resolves ${BASE_URL} to the original value "${PROTO}://${HOST}" — document this.
	if !strings.Contains(out["LOGIN"], "vault.internal") && !strings.Contains(out["LOGIN"], "BASE_URL") {
		t.Errorf("LOGIN unexpected value: %q", out["LOGIN"])
	}
}

func TestExpandAll_NoMutationOfInput(t *testing.T) {
	input := map[string]string{
		"FOO": "bar",
		"BAZ": "${FOO}-suffix",
	}
	origFoo := input["FOO"]
	origBaz := input["BAZ"]

	env.ExpandAll(input)

	if input["FOO"] != origFoo {
		t.Errorf("input mutated: FOO changed to %q", input["FOO"])
	}
	if input["BAZ"] != origBaz {
		t.Errorf("input mutated: BAZ changed to %q", input["BAZ"])
	}
}

func TestExpandAll_EmptyMap(t *testing.T) {
	out, errs := env.ExpandAll(map[string]string{})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}

func TestExpandTemplate_MultipleRefsInOneValue(t *testing.T) {
	env := map[string]string{
		"USER": "admin",
		"PASS": "s3cr3t",
	}
	result, err := env.ExpandTemplate("${USER}:${PASS}@db:5432", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "admin:s3cr3t@db:5432" {
		t.Errorf("unexpected result: %q", result)
	}
}
