package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// Change describes a single key-level difference between two snapshots.
type Change struct {
	Key    string
	OldVal string
	NewVal string
	Kind   ChangeKind
}

// ChangeKind classifies a snapshot difference.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// DiffEntries computes the differences between two snapshot entries.
func DiffEntries(old, new Entry) []Change {
	var changes []Change

	for k, nv := range new.Secrets {
		if ov, ok := old.Secrets[k]; !ok {
			changes = append(changes, Change{Key: k, NewVal: nv, Kind: Added})
		} else if ov != nv {
			changes = append(changes, Change{Key: k, OldVal: ov, NewVal: nv, Kind: Changed})
		}
	}
	for k, ov := range old.Secrets {
		if _, ok := new.Secrets[k]; !ok {
			changes = append(changes, Change{Key: k, OldVal: ov, Kind: Removed})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})
	return changes
}

// FormatDiff returns a human-readable summary of snapshot changes.
func FormatDiff(changes []Change) string {
	if len(changes) == 0 {
		return "no changes"
	}
	var sb strings.Builder
	for _, c := range changes {
		switch c.Kind {
		case Added:
			fmt.Fprintf(&sb, "+ %s\n", c.Key)
		case Removed:
			fmt.Fprintf(&sb, "- %s\n", c.Key)
		case Changed:
			fmt.Fprintf(&sb, "~ %s\n", c.Key)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
