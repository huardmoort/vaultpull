package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	Target    string    `json:"target"`
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
	file *os.File
}

// New opens (or creates) the audit log file at the given path.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{path: path, file: f}, nil
}

// Log writes a single entry as a JSON line.
func (l *Logger) Log(op, vaultPath, target, status, message string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Path:      vaultPath,
		Target:    target,
		Status:    status,
		Message:   message,
	}
	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.file, "%s\n", line)
	return err
}

// Close closes the underlying file.
func (l *Logger) Close() error {
	return l.file.Close()
}
