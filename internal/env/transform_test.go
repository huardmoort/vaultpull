package env

import (
	"testing"
)

func baseTransformSecrets() map[string]string {
	return map[string]string{
		"Api_Key":   "  abc123  ",
		"DB_Host":   "localhost",
		"db_pass":   " secret ",
	}
}

func TestTransform_UpperCase(t *testing.T) {
	tr := NewTransformer(CaseUpper, false)
	out := tr.Apply(baseTransformSecrets())

	if _, ok := out["API_KEY"]; !ok {
		t.Error("expected API_KEY in output")
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in output")
	}
	if _, ok := out["DB_PASS"]; !ok {
		t.Error("expected DB_PASS in output")
	}
}

func TestTransform_LowerCase(t *testing.T) {
	tr := NewTransformer(CaseLower, false)
	out := tr.Apply(baseTransformSecrets())

	if _, ok := out["api_key"]; !ok {
		t.Error("expected api_key in output")
	}
	if _, ok := out["db_host"]; !ok {
		t.Error("expected db_host in output")
	}
}

func TestTransform_NoCase(t *testing.T) {
	tr := NewTransformer(CaseNone, false)
	out := tr.Apply(baseTransformSecrets())

	if _, ok := out["Api_Key"]; !ok {
		t.Error("expected original key Api_Key preserved")
	}
}

func TestTransform_TrimValues(t *testing.T) {
	tr := NewTransformer(CaseNone, true)
	out := tr.Apply(baseTransformSecrets())

	if got := out["Api_Key"]; got != "abc123" {
		t.Errorf("expected trimmed value 'abc123', got %q", got)
	}
	if got := out["db_pass"]; got != "secret" {
		t.Errorf("expected trimmed value 'secret', got %q", got)
	}
}

func TestTransform_NoMutation(t *testing.T) {
	original := baseTransformSecrets()
	originalCopy := make(map[string]string, len(original))
	for k, v := range original {
		originalCopy[k] = v
	}

	tr := NewTransformer(CaseUpper, true)
	tr.Apply(original)

	for k, v := range originalCopy {
		if original[k] != v {
			t.Errorf("original map was mutated at key %q", k)
		}
	}
}

func TestTransform_EmptyMap(t *testing.T) {
	tr := NewTransformer(CaseUpper, true)
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
