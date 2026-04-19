package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Mapping defines a single Vault path -> local env file relationship.
type Mapping struct {
	VaultPath string `yaml:"vault_path"`
	EnvFile   string `yaml:"env_file"`
}

// Config holds the full vaultpull configuration.
type Config struct {
	VaultAddr  string    `yaml:"vault_addr"`
	VaultToken string    `yaml:"vault_token"`
	Mappings   []Mapping `yaml:"mappings"`
}

// Load reads and validates a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config: file not found: %s", path)
		}
		return nil, fmt.Errorf("config: read error: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse error: %w", err)
	}

	if cfg.VaultAddr == "" {
		return nil, errors.New("config: vault_addr is required")
	}
	if len(cfg.Mappings) == 0 {
		return nil, errors.New("config: at least one mapping is required")
	}
	for i, m := range cfg.Mappings {
		if m.VaultPath == "" {
			return nil, fmt.Errorf("config: mapping[%d]: vault_path is required", i)
		}
		if m.EnvFile == "" {
			return nil, fmt.Errorf("config: mapping[%d]: env_file is required", i)
		}
	}

	return &cfg, nil
}
