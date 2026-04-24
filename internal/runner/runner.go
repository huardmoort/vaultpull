package runner

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/lock"
	"github.com/yourusername/vaultpull/internal/vault"
	"github.com/yourusername/vaultpull/internal/writer"
)

// VaultReader is satisfied by *vault.Client.
type VaultReader interface {
	ReadSecret(path string) (map[string]string, error)
}

// Result holds the outcome of processing a single mapping.
type Result struct {
	EnvFile string
	Keys    int
	Err     error
}

// Run executes the full sync for all mappings defined in cfg.
func Run(cfg *config.Config) ([]Result, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}
	return runWithClient(cfg, client)
}

func runWithClient(cfg *config.Config, client VaultReader) ([]Result, error) {
	var results []Result
	for _, m := range cfg.Mappings {
		res := processMapping(client, m)
		results = append(results, res)
	}
	return results, nil
}

func processMapping(client VaultReader, m config.Mapping) Result {
	res := Result{EnvFile: m.EnvFile}

	lf, err := lock.Acquire(m.EnvFile)
	if err != nil {
		res.Err = fmt.Errorf("acquire lock: %w", err)
		return res
	}
	defer lf.Release()

	secrets, err := vault.FetchAll(client, m.Secrets)
	if err != nil {
		res.Err = fmt.Errorf("fetch secrets: %w", err)
		return res
	}

	if err := writer.MergeEnvFile(m.EnvFile, secrets); err != nil {
		res.Err = fmt.Errorf("write env file: %w", err)
		return res
	}

	res.Keys = len(secrets)
	return res
}
