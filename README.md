# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files safely

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Authenticate with Vault and pull secrets into a local `.env` file:

```bash
# Set your Vault address and token
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-token-here"

# Pull secrets from a Vault path into a .env file
vaultpull pull --path secret/data/myapp --out .env
```

**Example `.env` output:**

```
DATABASE_URL=postgres://user:pass@localhost/db
API_KEY=abc123
DEBUG=false
```

### Available Commands

| Command | Description |
|---------|-------------|
| `pull` | Sync secrets from Vault to a local `.env` file |
| `diff` | Preview changes before writing |
| `version` | Print the current version |

### Flags

```
--path    Vault secret path (required)
--out     Output file path (default: .env)
--dry-run Preview secrets without writing to disk
```

---

## Configuration

`vaultpull` respects standard Vault environment variables:

- `VAULT_ADDR` — Vault server address
- `VAULT_TOKEN` — Authentication token
- `VAULT_NAMESPACE` — Vault namespace (Enterprise)

---

## License

MIT © [yourusername](https://github.com/yourusername)