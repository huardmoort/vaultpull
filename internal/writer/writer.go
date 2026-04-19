package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteEnvFile writes key-value pairs to a .env file at the given path.
// Existing file content is replaced. The directory is created if needed.
func WriteEnvFile(path string, secrets map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", path, err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	for k, v := range secrets {
		line := fmt.Sprintf("%s=%s\n", k, escapeValue(v))
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("writing key %s to %s: %w", k, path, err)
		}
	}
	return nil
}

// escapeValue wraps the value in double quotes if it contains spaces or
// special characters that would break naive .env parsers.
func escapeValue(v string) string {
	if strings.ContainsAny(v, " \t\n\r#") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
