package env_test

import (
	"testing"

	"vaultpull/internal/env"
)

// TestFilter_ChainedIncludeExclude verifies that applying an include filter
// followed by an exclude filter yields the correct final key set.
func TestFilter_ChainedIncludeExclude(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":        "localhost",
		"DB_PASSWORD":    "s3cr3t",
		"DB_NAME":        "mydb",
		"APP_KEY":        "abc",
		"APP_SECRET_KEY": "xyz",
	}

	include := env.NewIncludeFilter([]string{"DB_*"})
	partial := include.Apply(secrets)

	exclude := env.NewExcludeFilter([]string{"*PASSWORD"})
	final := exclude.Apply(partial)

	if len(final) != 2 {
		t.Fatalf("expected 2 keys after chained filters, got %d: %v", len(final), final)
	}
	if _, ok := final["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in final result")
	}
	if _, ok := final["DB_NAME"]; !ok {
		t.Error("expected DB_NAME in final result")
	}
}

// TestFilter_WithRedactMap ensures filtered secrets can be safely redacted.
func TestFilter_WithRedactMap(t *testing.T) {
	secrets := map[string]string{
		"API_KEY":    "real-key-value",
		"API_SECRET": "real-secret",
		"LOG_LEVEL":  "debug",
	}

	f := env.NewIncludeFilter([]string{"API_*"})
	filtered := f.Apply(secrets)
	redacted := env.RedactMap(filtered)

	for k, v := range redacted {
		if env.IsSensitive(k) && v != "[REDACTED]" {
			t.Errorf("key %s should be redacted, got %q", k, v)
		}
	}
	if _, ok := redacted["LOG_LEVEL"]; ok {
		t.Error("LOG_LEVEL should not be present after include filter on API_*")
	}
}

// TestFilter_NoMutationOfInput confirms Apply does not modify the original map.
func TestFilter_NoMutationOfInput(t *testing.T) {
	original := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
	}
	copy := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
	}

	f := env.NewExcludeFilter([]string{"*PASS"})
	_ = f.Apply(original)

	for k, v := range copy {
		if original[k] != v {
			t.Errorf("original map was mutated at key %s", k)
		}
	}
	if len(original) != len(copy) {
		t.Error("original map length changed after Apply")
	}
}
