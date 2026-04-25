package env_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
)

// TestValidate_MultipleErrors ensures all bad keys are reported together.
func TestValidate_MultipleErrors(t *testing.T) {
	secrets := map[string]string{
		"bad-key-one": "v1",
		"123STARTS_NUM": "v2",
		"GOOD_ONE":     "v3",
	}
	_, err := env.Validate(secrets)
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve, ok := err.(*env.ValidationError)
	if !ok {
		t.Fatalf("wrong error type: %T", err)
	}
	if len(ve.Problems) != 2 {
		t.Errorf("expected 2 problems, got %d: %v", len(ve.Problems), ve.Problems)
	}
	if !strings.Contains(ve.Error(), "env validation failed") {
		t.Errorf("unexpected error message: %s", ve.Error())
	}
}

// TestValidate_MixedWarningsAndErrors checks that results are still returned
// alongside a hard error so callers can display partial context.
func TestValidate_MixedWarningsAndErrors(t *testing.T) {
	secrets := map[string]string{
		"bad-key":    "value",
		"EMPTY_GOOD": "",
	}
	results, err := env.Validate(secrets)
	if err == nil {
		t.Fatal("expected error")
	}
	// EMPTY_GOOD should still produce a result with a warning.
	var found bool
	for _, r := range results {
		if r.Key == "EMPTY_GOOD" && r.Warning != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning result for EMPTY_GOOD even when other keys are invalid")
	}
}

// TestFormatWarnings_Integration exercises FormatWarnings with realistic data.
func TestFormatWarnings_Integration(t *testing.T) {
	secrets := map[string]string{
		"TOKEN":        "",
		"DATABASE_URL": "postgres://localhost",
	}
	results, err := env.Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := env.FormatWarnings(results)
	if !strings.Contains(out, "TOKEN") {
		t.Errorf("expected TOKEN in warnings output, got:\n%s", out)
	}
	if strings.Contains(out, "DATABASE_URL") {
		t.Errorf("DATABASE_URL should not appear in warnings, got:\n%s", out)
	}
}
