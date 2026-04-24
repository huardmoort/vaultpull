package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const lockFileName = ".vaultpull.lock"

// LockFile represents an acquired lock on a target directory.
type LockFile struct {
	path string
}

// Acquire creates a lock file in the directory containing envPath.
// Returns an error if the lock is already held.
func Acquire(envPath string) (*LockFile, error) {
	dir := filepath.Dir(envPath)
	lockPath := filepath.Join(dir, lockFileName)

	if data, err := os.ReadFile(lockPath); err == nil {
		parts := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
		pid := parts[0]
		timestamp := ""
		if len(parts) > 1 {
			timestamp = parts[1]
		}
		return nil, fmt.Errorf("lock already held by pid %s (acquired %s); remove %s to force", pid, timestamp, lockPath)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create lock dir: %w", err)
	}

	content := fmt.Sprintf("%d\n%s", os.Getpid(), time.Now().UTC().Format(time.RFC3339))
	if err := os.WriteFile(lockPath, []byte(content), 0o644); err != nil {
		return nil, fmt.Errorf("write lock file: %w", err)
	}

	return &LockFile{path: lockPath}, nil
}

// Release removes the lock file.
func (l *LockFile) Release() error {
	if err := os.Remove(l.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("release lock: %w", err)
	}
	return nil
}

// IsHeld returns true if a lock file exists for the given envPath.
func IsHeld(envPath string) bool {
	dir := filepath.Dir(envPath)
	lockPath := filepath.Join(dir, lockFileName)
	_, err := os.Stat(lockPath)
	return err == nil
}

// OwnerPID returns the PID stored in the lock file, or 0 if not held.
func OwnerPID(envPath string) int {
	dir := filepath.Dir(envPath)
	lockPath := filepath.Join(dir, lockFileName)
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return 0
	}
	parts := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
	pid, _ := strconv.Atoi(parts[0])
	return pid
}
