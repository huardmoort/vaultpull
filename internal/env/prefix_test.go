package env

import (
	"testing"
)

func TestAddPrefix_Basic(t *testing.T) {
	tr := NewPrefixTransformer("APP_")
	input := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	got := tr.AddPrefix(input)

	if got["APP_DB_HOST"] != "localhost" {
		t.Errorf("expected APP_DB_HOST=localhost, got %q", got["APP_DB_HOST"])
	}
	if got["APP_DB_PORT"] != "5432" {
		t.Errorf("expected APP_DB_PORT=5432, got %q", got["APP_DB_PORT"])
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestAddPrefix_AlreadyPrefixed(t *testing.T) {
	tr := NewPrefixTransformer("APP_")
	input := map[string]string{"APP_TOKEN": "secret"}
	got := tr.AddPrefix(input)

	if _, ok := got["APP_TOKEN"]; !ok {
		t.Error("expected APP_TOKEN to be present without double-prefix")
	}
	if _, ok := got["APP_APP_TOKEN"]; ok {
		t.Error("did not expect APP_APP_TOKEN (double prefix)")
	}
}

func TestStripPrefix_Basic(t *testing.T) {
	tr := NewPrefixTransformer("APP_")
	input := map[string]string{"APP_DB_HOST": "localhost", "OTHER_KEY": "value"}
	got := tr.StripPrefix(input)

	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if _, ok := got["OTHER_KEY"]; ok {
		t.Error("expected OTHER_KEY to be dropped (no matching prefix)")
	}
	if len(got) != 1 {
		t.Errorf("expected 1 key, got %d", len(got))
	}
}

func TestStripPrefix_NoMatches(t *testing.T) {
	tr := NewPrefixTransformer("APP_")
	input := map[string]string{"FOO": "bar", "BAZ": "qux"}
	got := tr.StripPrefix(input)

	if len(got) != 0 {
		t.Errorf("expected empty map, got %d keys", len(got))
	}
}

func TestRenamePrefix(t *testing.T) {
	input := map[string]string{
		"OLD_DB_HOST": "localhost",
		"OLD_DB_PORT": "5432",
		"UNRELATED":   "value",
	}
	got := RenamePrefix(input, "OLD_", "NEW_")

	if got["NEW_DB_HOST"] != "localhost" {
		t.Errorf("expected NEW_DB_HOST=localhost, got %q", got["NEW_DB_HOST"])
	}
	if got["NEW_DB_PORT"] != "5432" {
		t.Errorf("expected NEW_DB_PORT=5432, got %q", got["NEW_DB_PORT"])
	}
	if got["UNRELATED"] != "value" {
		t.Errorf("expected UNRELATED=value, got %q", got["UNRELATED"])
	}
	if _, ok := got["OLD_DB_HOST"]; ok {
		t.Error("did not expect OLD_DB_HOST to remain")
	}
}

func TestRenamePrefix_EmptyOld(t *testing.T) {
	input := map[string]string{"FOO": "bar"}
	got := RenamePrefix(input, "", "NEW_")

	if got["NEW_FOO"] != "bar" {
		t.Errorf("expected NEW_FOO=bar, got %q", got["NEW_FOO"])
	}
}
