package scripts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	hookScriptFilename         = "claude-hook.sh"
	notificationHookFilename   = "notification-hook.sh"
	plistFilename              = "com.dee.claude-notifier.plist"
	binaryPathMarker           = "__BINARY_PATH__"
)

// WriteHookScript writes claude-hook.sh to destDir with executable permissions (0755).
// It creates parent directories if they don't exist.
// Returns the full path to the written file.
func WriteHookScript(destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("create directory: %w", err)
	}

	destPath := filepath.Join(destDir, hookScriptFilename)

	if err := os.WriteFile(destPath, HookScript, 0755); err != nil {
		return "", fmt.Errorf("write hook script: %w", err)
	}

	return destPath, nil
}

// WriteNotificationHookScript writes notification-hook.sh to destDir with executable permissions (0755).
// It creates parent directories if they don't exist.
// Returns the full path to the written file.
func WriteNotificationHookScript(destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("create directory: %w", err)
	}

	destPath := filepath.Join(destDir, notificationHookFilename)

	if err := os.WriteFile(destPath, NotificationHookScript, 0755); err != nil {
		return "", fmt.Errorf("write notification hook script: %w", err)
	}

	return destPath, nil
}

// WritePlist writes the plist to destDir, substituting __BINARY_PATH__ with binaryPath.
// It creates parent directories if they don't exist.
// Returns the full path to the written file.
func WritePlist(destDir string, binaryPath string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("create directory: %w", err)
	}

	content := strings.ReplaceAll(string(PlistTemplate), binaryPathMarker, binaryPath)

	destPath := filepath.Join(destDir, plistFilename)

	if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write plist: %w", err)
	}

	return destPath, nil
}
