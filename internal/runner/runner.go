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
		secrets, fetchErr := vault.FetchAll(client, m.VaultPath)
		status := "ok"
		msg := ""
		if fetchErr != nil {
			status = "error"
			msg = fetchErr.Error()
			log.Printf("[warn] fetch %s: %v", m.VaultPath, fetchErr)
		} else {
			if writeErr := writer.MergeEnvFile(m.EnvFile, secrets); writeErr != nil {
				status = "error"
				msg = writeErr.Error()
				log.Printf("[warn] write %s: %v", m.EnvFile, writeErr)
			}
		}
		if logger != nil {
			if logErr := logger.Log("sync", m.VaultPath, m.EnvFile, status, msg); logErr != nil {
				log.Printf("[warn] audit log: %v", logErr)
			}
		}
	}
	return nil
}
