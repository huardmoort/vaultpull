package lock

import (
	"os"
	"path/filepath"
	"testing"
)

func tempEnvPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, ".env")
}

func TestAcquire_CreatesLockFile(t *testing.T) {
	envPath := tempEnvPath(t)
	lf, err := Acquire(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer lf.Release()

	if !IsHeld(envPath) {
		t.Error("expected lock to be held after Acquire")
	}
}

func TestAcquire_FailsWhenAlreadyLocked(t *testing.T) {
	envPath := tempEnvPath(t)
	lf, err := Acquire(envPath)
	if err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer lf.Release()

	_, err = Acquire(envPath)
	if err == nil {
		t.Fatal("expected error on second Acquire, got nil")
	}
}

func TestRelease_RemovesLockFile(t *testing.T) {
	envPath := tempEnvPath(t)
	lf, err := Acquire(envPath)
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	if err := lf.Release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}

	if IsHeld(envPath) {
		t.Error("expected lock to be released")
	}
}

func TestRelease_IdempotentOnMissingFile(t *testing.T) {
	envPath := tempEnvPath(t)
	lf, err := Acquire(envPath)
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	// Manually remove before Release
	os.Remove(filepath.Join(filepath.Dir(envPath), lockFileName))

	if err := lf.Release(); err != nil {
		t.Errorf("expected no error on missing lock file, got: %v", err)
	}
}

func TestOwnerPID_ReturnsCurrentPID(t *testing.T) {
	envPath := tempEnvPath(t)
	lf, err := Acquire(envPath)
	if err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	defer lf.Release()

	pid := OwnerPID(envPath)
	if pid != os.Getpid() {
		t.Errorf("expected pid %d, got %d", os.Getpid(), pid)
	}
}

func TestOwnerPID_ZeroWhenNotHeld(t *testing.T) {
	envPath := tempEnvPath(t)
	if pid := OwnerPID(envPath); pid != 0 {
		t.Errorf("expected 0, got %d", pid)
	}
}

func TestIsHeld_FalseWhenNoLock(t *testing.T) {
	envPath := tempEnvPath(t)
	if IsHeld(envPath) {
		t.Error("expected IsHeld to return false for fresh dir")
	}
}
