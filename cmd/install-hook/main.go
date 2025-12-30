package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeroendee/claude-notifier/internal/scripts"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	destDir := filepath.Join(homeDir, ".claude", "hooks")
	if err := run(destDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(destDir string) error {
	path, err := scripts.WriteHookScript(destDir)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote hook script to: %s\n", path)
	return nil
}
