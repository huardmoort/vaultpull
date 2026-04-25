package env

import (
	"strings"
	"testing"
)

func TestValidate_AllClean(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"API_KEY":      "abc123",
	}
	results, err := Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Warning != "" {
			t.Errorf("unexpected warning for %s: %s", r.Key, r.Warning)
		}
	}
}

func TestValidate_InvalidKeyName(t *testing.T) {
	secrets := map[string]string{
		"bad-key":  "value",
		"GOOD_KEY": "value",
	}
	_, err := Validate(secrets)
	if err == nil {
		t.Fatal("expected error for invalid key name")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Problems) != 1 {
		t.Errorf("expected 1 problem, got %d: %v", len(ve.Problems), ve.Problems)
	}
}

func TestValidate_EmptyValueWarning(t *testing.T) {
	secrets := map[string]string{
		"EMPTY_VAR": "",
	}
	results, err := Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Warning == "" {
		t.Error("expected a warning for empty value")
	}
}

func TestValidate_NewlineWarning(t *testing.T) {
	secrets := map[string]string{
		"MULTI_LINE": "line1\nline2",
	}
	results, err := Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 || !strings.Contains(results[0].Warning, "newline") {
		t.Error("expected newline warning")
	}
}

func TestValidate_LargeValueWarning(t *testing.T) {
	secrets := map[string]string{
		"BIG_SECRET": strings.Repeat("x", 5000),
	}
	results, err := Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 || !strings.Contains(results[0].Warning, "unusually large") {
		t.Error("expected large-value warning")
	}
}

func TestFormatWarnings_NoWarnings(t *testing.T) {
	results := []Result{{Key: "A", Warning: ""}, {Key: "B", Warning: ""}}
	out := FormatWarnings(results)
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}

func TestFormatWarnings_WithWarning(t *testing.T) {
	results := []Result{{Key: "EMPTY_VAR", Warning: "value is empty"}}
	out := FormatWarnings(results)
	if !strings.Contains(out, "EMPTY_VAR") || !strings.Contains(out, "value is empty") {
		t.Errorf("unexpected format output: %q", out)
	}
}
