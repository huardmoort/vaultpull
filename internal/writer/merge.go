package writer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MergeEnvFile merges secrets into an existing .env file, preserving comments
// and unknown keys, while overwriting keys present in secrets.
func MergeEnvFile(path string, secrets map[string]string) error {
	existing, lines, err := readEnvFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading existing env file: %w", err)
	}

	// Track which secret keys have been written via update
	updated := make(map[string]bool)

	var out []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			out = append(out, line)
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			out = append(out, line)
			continue
		}
		key := parts[0]
		if val, ok := secrets[key]; ok {
			out = append(out, fmt.Sprintf("%s=%s", key, escapeValue(val)))
			updated[key] = true
		} else {
			out = append(out, line)
		}
	}

	// Append new keys not present in existing file
	for key, val := range secrets {
		if !updated[key] {
			if _, exists := existing[key]; !exists {
				out = append(out, fmt.Sprintf("%s=%s", key, escapeValue(val)))
			}
		}
	}

	content := strings.Join(out, "\n")
	if len(out) > 0 {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0600)
}

// readEnvFile reads an env file and returns a map of key->value and the raw lines.
func readEnvFile(path string) (map[string]string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	kvs := make(map[string]string)
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			kvs[parts[0]] = parts[1]
		}
	}
	return kvs, lines, scanner.Err()
}
