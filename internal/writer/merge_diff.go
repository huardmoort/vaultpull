package writer

import (
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/diff"
)

// MergeWithDiff merges new secrets into the env file and prints a diff summary
// to stdout. It returns the list of changes applied.
func MergeWithDiff(path string, incoming map[string]string, verbose bool) ([]diff.Change, error) {
	existing, err := readEnvFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading existing env file: %w", err)
	}

	changes := diff.Compute(existing, incoming)

	if err := MergeEnvFile(path, incoming); err != nil {
		return nil, err
	}

	if verbose {
		fmt.Printf("[vaultpull] %s:\n%s", path, diff.Format(changes))
	}

	return changes, nil
}
