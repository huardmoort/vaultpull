package env

import (
	"fmt"
	"regexp"
	"strings"
)

// validKeyRe matches legal environment variable names: [A-Z_][A-Z0-9_]*
var validKeyRe = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// ValidationError holds all problems found during validation.
type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("env validation failed:\n  %s", strings.Join(e.Problems, "\n  "))
}

// Result holds the outcome of validating a single key/value pair.
type Result struct {
	Key     string
	Warning string // non-empty when the value looks suspicious but is not fatal
}

// Validate checks a map of env key→value pairs for common problems.
// It returns a ValidationError if any hard errors are found, plus a slice
// of Results that carry non-fatal warnings.
func Validate(secrets map[string]string) ([]Result, error) {
	var problems []string
	var results []Result

	for k, v := range secrets {
		if !validKeyRe.MatchString(k) {
			problems = append(problems, fmt.Sprintf("invalid key name %q (must match [A-Z_][A-Z0-9_]*)", k))
			continue
		}

		r := Result{Key: k}

		switch {
		case v == "":
			r.Warning = "value is empty"
		case strings.Contains(v, "\n"):
			r.Warning = "value contains newline characters"
		case len(v) > 4096:
			r.Warning = fmt.Sprintf("value is unusually large (%d bytes)", len(v))
		}

		results = append(results, r)
	}

	if len(problems) > 0 {
		return results, &ValidationError{Problems: problems}
	}
	return results, nil
}

// FormatWarnings returns a human-readable summary of any warnings in results.
func FormatWarnings(results []Result) string {
	var sb strings.Builder
	for _, r := range results {
		if r.Warning != "" {
			fmt.Fprintf(&sb, "  WARN  %s: %s\n", r.Key, r.Warning)
		}
	}
	return sb.String()
}
