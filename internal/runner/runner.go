package runner

import (
	"fmt"
	"log"

	"github.com/yourorg/vaultpull/internal/audit"
	"github.com/yourorg/vaultpull/internal/config"
	"github.com/yourorg/vaultpull/internal/vault"
	"github.com/yourorg/vaultpull/internal/writer"
)

// Run executes the full sync: fetch secrets from Vault and write .env files.
// It processes all mappings defined in cfg, logging results to auditLog if non-empty.
// Errors for individual mappings are logged as warnings; Run itself returns nil
// unless a fatal setup error occurs (e.g. Vault client or audit logger init).
func Run(cfg *config.Config, auditLog string) error {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("runner: vault client: %w", err)
	}

	var logger *audit.Logger
	if auditLog != "" {
		logger, err = audit.New(auditLog)
		if err != nil {
			return fmt.Errorf("runner: audit logger: %w", err)
		}
		defer logger.Close()
	}

	for _, m := range cfg.Mappings {
		processMapping(client, logger, m)
	}
	return nil
}

// processMapping fetches secrets for a single mapping and writes them to the
// corresponding env file, logging the outcome to the audit logger if set.
func processMapping(client *vault.Client, logger *audit.Logger, m config.Mapping) {
	status := "ok"
	msg := ""

	secrets, fetchErr := vault.FetchAll(client, m.VaultPath)
	if fetchErr != nil {
		status = "error"
		msg = fetchErr.Error()
		log.Printf("[warn] fetch %s: %v", m.VaultPath, fetchErr)
	} else if writeErr := writer.MergeEnvFile(m.EnvFile, secrets); writeErr != nil {
		status = "error"
		msg = writeErr.Error()
		log.Printf("[warn] write %s: %v", m.EnvFile, writeErr)
	}

	if logger != nil {
		if logErr := logger.Log("sync", m.VaultPath, m.EnvFile, status, msg); logErr != nil {
			log.Printf("[warn] audit log: %v", logErr)
		}
	}
}
