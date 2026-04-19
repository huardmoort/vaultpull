package diff

import (
	"strings"
	"testing"
)

func TestCompute_Add(t *testing.T) {
	old := map[string]string{}
	new_ := map[string]string{"FOO": "bar"}
	changes := Compute(old, new_)
	if len(changes) != 1 || changes[0].Action != "add" || changes[0].Key != "FOO" {
		t.Fatalf("expected one add change, got %+v", changes)
	}
}

func TestCompute_Update(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new_ := map[string]string{"FOO": "baz"}
	changes := Compute(old, new_)
	if len(changes) != 1 || changes[0].Action != "update" {
		t.Fatalf("expected one update change, got %+v", changes)
	}
	if changes[0].Old != "bar" || changes[0].New != "baz" {
		t.Fatalf("unexpected old/new values: %+v", changes[0])
	}
}

func TestCompute_Delete(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new_ := map[string]string{}
	changes := Compute(old, new_)
	if len(changes) != 1 || changes[0].Action != "delete" {
		t.Fatalf("expected one delete change, got %+v", changes)
	}
}

func TestCompute_NoChange(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new_ := map[string]string{"FOO": "bar"}
	changes := Compute(old, new_)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %+v", changes)
	}
}

func TestFormat_Empty(t *testing.T) {
	out := Format(nil)
	if out != "No changes." {
		t.Fatalf("expected 'No changes.', got %q", out)
	}
}

func TestFormat_Mixed(t *testing.T) {
	changes := []Change{
		{Key: "A", Action: "add"},
		{Key: "B", Action: "update"},
		{Key: "C", Action: "delete"},
	}
	out := Format(changes)
	if !strings.Contains(out, "+ A") {
		t.Error("expected '+ A'")
	}
	if !strings.Contains(out, "~ B") {
		t.Error("expected '~ B'")
	}
	if !strings.Contains(out, "- C") {
		t.Error("expected '- C'")
	}
}
