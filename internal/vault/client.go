package vault

import (
	"fmt"
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	logical *vaultapi.Logical
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(addr, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	cfg.HttpClient = &http.Client{Timeout: 10 * time.Second}

	c, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to create client: %w", err)
	}
	c.SetToken(token)

	return &Client{logical: c.Logical()}, nil
}

// ReadSecret reads a KV secret at the given path and returns its data map.
func (c *Client) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: no secret found at %q", path)
	}

	// KV v2 wraps data under a "data" key.
	if data, ok := secret.Data["data"]; ok {
		if m, ok := data.(map[string]interface{}); ok {
			return m, nil
		}
	}

	return secret.Data, nil
}
