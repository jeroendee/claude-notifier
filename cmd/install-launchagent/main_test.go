package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_Success(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()
	binaryPath := "/usr/local/bin/claude-notifier"

	err := run(destDir, binaryPath)
	if err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}

	plistPath := filepath.Join(destDir, "com.dee.claude-notifier.plist")
	content, err := os.ReadFile(plistPath)
	if err != nil {
		t.Fatalf("read plist: %v", err)
	}

	if !strings.Contains(string(content), binaryPath) {
		t.Error("plist does not contain binary path")
	}

	if strings.Contains(string(content), "__BINARY_PATH__") {
		t.Error("plist still contains placeholder")
	}
}

func TestRun_InvalidDir(t *testing.T) {
	t.Parallel()

	err := run("/dev/null/launchagents", "/usr/local/bin/test")
	if err == nil {
		t.Error("run() error = nil, want error for invalid directory")
	}
}
