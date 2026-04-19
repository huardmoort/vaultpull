package vault

import (
	"fmt"

	"github.com/user/vaultpull/internal/config"
)

// SecretMap holds key-value pairs resolved from Vault.
type SecretMap map[string]string

// FetchAll reads all secrets defined in the config mappings and returns
// a map of env-file destination path -> SecretMap.
func FetchAll(c *Client, mappings []config.Mapping) (map[string]SecretMap, error) {
	result := make(map[string]SecretMap, len(mappings))

	for _, m := range mappings {
		data, err := c.ReadSecret(m.VaultPath)
		if err != nil {
			return nil, fmt.Errorf("fetch: %w", err)
		}

		sm := make(SecretMap, len(data))
		for k, v := range data {
			str, ok := v.(string)
			if !ok {
				str = fmt.Sprintf("%v", v)
			}
			sm[k] = str
		}
		result[m.EnvFile] = sm
	}

	return result, nil
}
