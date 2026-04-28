package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// templateVarRe matches ${VAR_NAME} or $VAR_NAME style references.
var templateVarRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// ExpandTemplate replaces template variable references in a string value
// using the provided env map, falling back to OS environment variables.
// Returns an error if any referenced variable is unresolved.
func ExpandTemplate(value string, env map[string]string) (string, error) {
	var missing []string

	result := templateVarRe.ReplaceAllStringFunc(value, func(match string) string {
		name := extractVarName(match)
		if v, ok := env[name]; ok {
			return v
		}
		if v, ok := os.LookupEnv(name); ok {
			return v
		}
		missing = append(missing, name)
		return match
	})

	if len(missing) > 0 {
		return "", fmt.Errorf("unresolved template variables: %s", strings.Join(missing, ", "))
	}
	return result, nil
}

// ExpandAll applies ExpandTemplate to every value in the map.
// It returns the expanded map and a slice of per-key errors.
func ExpandAll(env map[string]string) (map[string]string, []error) {
	out := make(map[string]string, len(env))
	var errs []error
	for k, v := range env {
		expanded, err := ExpandTemplate(v, env)
		if err != nil {
			errs = append(errs, fmt.Errorf("key %q: %w", k, err))
			out[k] = v // preserve original on error
		} else {
			out[k] = expanded
		}
	}
	return out, errs
}

func extractVarName(match string) string {
	// Strip ${ and } or just $
	match = strings.TrimPrefix(match, "${")
	match = strings.TrimSuffix(match, "}")
	match = strings.TrimPrefix(match, "$")
	return match
}
