package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jeroendee/claude-notifier/internal/scripts"
)

const binaryName = "claude-notifier"

func main() {
	binaryPath, err := findBinaryPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding binary: %v\n", err)
		os.Exit(1)
	}

	destDir := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")
	if err := run(destDir, binaryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(destDir, binaryPath string) error {
	path, err := scripts.WritePlist(destDir, binaryPath)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote plist to: %s\n", path)
	return nil
}

func findBinaryPath() (string, error) {
	gobin := os.Getenv("GOBIN")
	if gobin == "" {
		out, err := exec.Command("go", "env", "GOPATH").Output()
		if err != nil {
			return "", fmt.Errorf("get GOPATH: %w", err)
		}
		gobin = filepath.Join(strings.TrimSpace(string(out)), "bin")
	}
	return filepath.Join(gobin, binaryName), nil
}
