package runner

import (
	"fmt"

	"github.com/vaultpull/internal/config"
	"github.com/vaultpull/internal/vault"
	"github.com/vaultpull/internal/writer"
)

// Result holds the outcome of syncing a single mapping.
type Result struct {
	EnvFile string
	Written int
	Err     error
}

// Run executes the full sync: fetch secrets from Vault and write .env files
// for every mapping defined in cfg.
func Run(cfg *config.Config) []Result {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return []Result{{Err: fmt.Errorf("vault client: %w", err)}}
	}

	results := make([]Result, 0, len(cfg.Mappings))

	for _, m := range cfg.Mappings {
		r := Result{EnvFile: m.EnvFile}

		secrets, err := vault.FetchAll(client, m.Secrets)
		if err != nil {
			r.Err = fmt.Errorf("fetch %s: %w", m.EnvFile, err)
			results = append(results, r)
			continue
		}

		if err := writer.MergeEnvFile(m.EnvFile, secrets); err != nil {
			r.Err = fmt.Errorf("write %s: %w", m.EnvFile, err)
			results = append(results, r)
			continue
		}

		r.Written = len(secrets)
		results = append(results, r)
	}

	return results
}
