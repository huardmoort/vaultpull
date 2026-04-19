package writer

import (
	"bufio"
	"os"
	"strings"
)

// MergeEnvFile reads an existing .env file and merges new secrets into it.
// Keys present in secrets overwrite existing values; other lines are preserved.
// If the file does not exist it is created from scratch.
func MergeEnvFile(path string, secrets map[string]string) error {
	existing, err := readEnvFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for k, v := range secrets {
		existing[k] = v
	}

	return WriteEnvFile(path, existing)
}

// readEnvFile parses a .env file into a map. Lines starting with '#' and
// blank lines are ignored. Only KEY=VALUE format is supported.
func readEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return result, scanner.Err()
}
