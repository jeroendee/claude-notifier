package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_Success(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()

	err := run(destDir)
	if err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}

	scriptPath := filepath.Join(destDir, "claude-hook.sh")
	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("stat hook script: %v", err)
	}

	if info.Mode().Perm() != 0755 {
		t.Errorf("permissions = %o, want 0755", info.Mode().Perm())
	}
}

func TestRun_InvalidDir(t *testing.T) {
	t.Parallel()

	// /dev/null is a file, not a directory - cannot create subdirectories
	err := run("/dev/null/hooks")
	if err == nil {
		t.Error("run() error = nil, want error for invalid directory")
	}
}
