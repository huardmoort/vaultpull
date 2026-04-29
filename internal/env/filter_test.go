package env

import (
	"testing"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"APP_KEY":     "abc123",
		"APP_SECRET":  "xyz",
		"LOG_LEVEL":   "info",
	}
}

func TestFilter_IncludeExactMatch(t *testing.T) {
	f := NewIncludeFilter([]string{"DB_HOST", "LOG_LEVEL"})
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in output")
	}
	if _, ok := out["LOG_LEVEL"]; !ok {
		t.Error("expected LOG_LEVEL in output")
	}
}

func TestFilter_IncludePrefixWildcard(t *testing.T) {
	f := NewIncludeFilter([]string{"DB_*"})
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(out), out)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST")
	}
	if _, ok := out["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD")
	}
}

func TestFilter_IncludeSuffixWildcard(t *testing.T) {
	f := NewIncludeFilter([]string{"*_SECRET"})
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(out), out)
	}
}

func TestFilter_IncludeContainsWildcard(t *testing.T) {
	f := NewIncludeFilter([]string{"*APP*"})
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Fatalf("expected 2 keys (APP_KEY, APP_SECRET), got %d", len(out))
	}
}

func TestFilter_ExcludePattern(t *testing.T) {
	f := NewExcludeFilter([]string{"*PASSWORD", "*SECRET"})
	out := f.Apply(baseSecrets())
	for k := range out {
		if k == "DB_PASSWORD" || k == "APP_SECRET" {
			t.Errorf("key %s should have been excluded", k)
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 remaining keys, got %d", len(out))
	}
}

func TestFilter_EmptyPatterns_Include(t *testing.T) {
	f := NewIncludeFilter([]string{})
	out := f.Apply(baseSecrets())
	if len(out) != 0 {
		t.Errorf("expected empty map with no include patterns, got %d", len(out))
	}
}

func TestFilter_EmptyPatterns_Exclude(t *testing.T) {
	f := NewExcludeFilter([]string{})
	out := f.Apply(baseSecrets())
	if len(out) != len(baseSecrets()) {
		t.Errorf("expected all keys preserved, got %d", len(out))
	}
}

func TestFilter_CaseInsensitiveMatch(t *testing.T) {
	f := NewIncludeFilter([]string{"db_*"})
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Errorf("expected case-insensitive match to find 2 keys, got %d", len(out))
	}
}
