package env

import (
	"os"
	"testing"
)

func TestExpandTemplate_NoVars(t *testing.T) {
	result, err := ExpandTemplate("hello world", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("expected 'hello world', got %q", result)
	}
}

func TestExpandTemplate_BraceStyle(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	result, err := ExpandTemplate("postgres://${HOST}:${PORT}/db", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "postgres://localhost:5432/db" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestExpandTemplate_BareStyle(t *testing.T) {
	env := map[string]string{"REGION": "us-east-1"}
	result, err := ExpandTemplate("aws/$REGION/secret", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "aws/us-east-1/secret" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestExpandTemplate_FallsBackToOS(t *testing.T) {
	os.Setenv("VAULTPULL_TEST_VAR", "from-os")
	defer os.Unsetenv("VAULTPULL_TEST_VAR")

	result, err := ExpandTemplate("${VAULTPULL_TEST_VAR}", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "from-os" {
		t.Errorf("expected 'from-os', got %q", result)
	}
}

func TestExpandTemplate_UnresolvedVar(t *testing.T) {
	_, err := ExpandTemplate("${MISSING_VAR}", map[string]string{})
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestExpandAll_Mixed(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api/v1",
		"PLAIN":    "no-vars-here",
	}
	out, errs := ExpandAll(env)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if out["API_URL"] != "https://example.com/api/v1" {
		t.Errorf("API_URL not expanded correctly: %q", out["API_URL"])
	}
	if out["PLAIN"] != "no-vars-here" {
		t.Errorf("PLAIN changed unexpectedly: %q", out["PLAIN"])
	}
}

func TestExpandAll_ReturnsOriginalOnError(t *testing.T) {
	env := map[string]string{
		"BAD": "${DOES_NOT_EXIST}",
	}
	out, errs := ExpandAll(env)
	if len(errs) == 0 {
		t.Fatal("expected errors")
	}
	if out["BAD"] != "${DOES_NOT_EXIST}" {
		t.Errorf("expected original value preserved, got %q", out["BAD"])
	}
}
