package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log record for a vaultpull run.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	ConfigPath string   `json:"config_path"`
	Added     int       `json:"added"`
	Updated   int       `json:"updated"`
	Deleted   int       `json:"deleted"`
	Error     string    `json:"error,omitempty"`
}

// Logger appends audit entries to a file.
type Logger struct {
	path string
}

// New creates a new Logger that writes to the given file path.
// The parent directory is created if it does not exist.
func New(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: create dir: %w", err)
	}
	return &Logger{path: path}, nil
}

// Log appends an Entry to the audit log file as a newline-delimited JSON record.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}
