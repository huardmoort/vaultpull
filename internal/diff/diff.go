package diff

import "fmt"

// Change represents a single key-level change between old and new env state.
type Change struct {
	Key    string
	Old    string
	New    string
	Action string // "add", "update", "delete"
}

// Compute returns the list of changes going from oldEnv to newEnv.
func Compute(oldEnv, newEnv map[string]string) []Change {
	var changes []Change

	for k, newVal := range newEnv {
		oldVal, exists := oldEnv[k]
		if !exists {
			changes = append(changes, Change{Key: k, Old: "", New: newVal, Action: "add"})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Old: oldVal, New: newVal, Action: "update"})
		}
	}

	for k, oldVal := range oldEnv {
		if _, exists := newEnv[k]; !exists {
			changes = append(changes, Change{Key: k, Old: oldVal, New: "", Action: "delete"})
		}
	}

	return changes
}

// Format returns a human-readable summary of changes.
func Format(changes []Change) string {
	if len(changes) == 0 {
		return "No changes."
	}
	var out string
	for _, c := range changes {
		switch c.Action {
		case "add":
			out += fmt.Sprintf("  + %s\n", c.Key)
		case "update":
			out += fmt.Sprintf("  ~ %s\n", c.Key)
		case "delete":
			out += fmt.Sprintf("  - %s\n", c.Key)
		}
	}
	return out
}
