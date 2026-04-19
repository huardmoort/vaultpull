package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the vaultpull configuration.
type Config struct {
	VaultAddr  string            `yaml:"vault_addr"`
	VaultToken string            `yaml:"vault_token"`
	Mappings   []SecretMapping   `yaml:"mappings"`
}

// SecretMapping maps a Vault secret path to a local .env file.
type SecretMapping struct {
	VaultPath string `yaml:"vault_path"`
	EnvFile   string `yaml:"env_file"`
}

// Load reads and parses the config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return fmt.Errorf("vault_addr is required")
	}
	if c.VaultToken == "" {
		return fmt.Errorf("vault_token is required")
	}
	if len(c.Mappings) == 0 {
		return fmt.Errorf("at least one mapping is required")
	}
	for i, m := range c.Mappings {
		if m.VaultPath == "" {
			return fmt.Errorf("mapping[%d]: vault_path is required", i)
		}
		if m.EnvFile == "" {
			return fmt.Errorf("mapping[%d]: env_file is required", i)
		}
	}
	return nil
}
