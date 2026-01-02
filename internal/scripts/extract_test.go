package scripts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteHookScript_CreatesFileWithCorrectContent(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()

	path, err := WriteHookScript(destDir)
	if err != nil {
		t.Fatalf("WriteHookScript() error = %v", err)
	}

	wantPath := filepath.Join(destDir, "claude-hook.sh")
	if path != wantPath {
		t.Errorf("WriteHookScript() path = %v, want %v", path, wantPath)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	if string(content) != string(HookScript) {
		t.Errorf("WriteHookScript() content mismatch")
	}
}

func TestWriteHookScript_SetsExecutablePermissions(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()

	path, err := WriteHookScript(destDir)
	if err != nil {
		t.Fatalf("WriteHookScript() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	wantMode := os.FileMode(0755)
	if info.Mode().Perm() != wantMode {
		t.Errorf("WriteHookScript() file mode = %v, want %v", info.Mode().Perm(), wantMode)
	}
}

func TestWriteHookScript_CreatesParentDirectories(t *testing.T) {
	t.Parallel()

	baseDir := t.TempDir()
	destDir := filepath.Join(baseDir, "nested", "dirs")

	path, err := WriteHookScript(destDir)
	if err != nil {
		t.Fatalf("WriteHookScript() error = %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("WriteHookScript() file does not exist at %v", path)
	}
}

func TestWriteHookScript_ReturnsErrorOnMkdirAllFailure(t *testing.T) {
	t.Parallel()

	// Use a non-existent parent with a file blocker
	baseDir := t.TempDir()
	blocker := filepath.Join(baseDir, "blocker")
	if err := os.WriteFile(blocker, []byte{}, 0644); err != nil {
		t.Fatalf("failed to create blocker file: %v", err)
	}
	// Try to write into a path where a file exists (can't create dir)
	destDir := filepath.Join(blocker, "subdir")

	_, err := WriteHookScript(destDir)
	if err == nil {
		t.Error("WriteHookScript() expected error for invalid path, got nil")
	}
}

func TestWriteHookScript_ReturnsErrorOnWriteFileFailure(t *testing.T) {
	t.Parallel()

	// Create a directory where the file should be written.
	// WriteFile will fail because it cannot write to a directory.
	destDir := t.TempDir()
	targetPath := filepath.Join(destDir, "claude-hook.sh")
	if err := os.Mkdir(targetPath, 0755); err != nil {
		t.Fatalf("failed to create blocking directory: %v", err)
	}

	_, err := WriteHookScript(destDir)
	if err == nil {
		t.Error("WriteHookScript() expected error when target is a directory, got nil")
	}
}

func TestWriteNotificationHookScript_CreatesFileWithCorrectContent(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()

	path, err := WriteNotificationHookScript(destDir)
	if err != nil {
		t.Fatalf("WriteNotificationHookScript() error = %v", err)
	}

	wantPath := filepath.Join(destDir, "notification-hook.sh")
	if path != wantPath {
		t.Errorf("WriteNotificationHookScript() path = %v, want %v", path, wantPath)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	if string(content) != string(NotificationHookScript) {
		t.Errorf("WriteNotificationHookScript() content mismatch")
	}
}

func TestWriteNotificationHookScript_SetsExecutablePermissions(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()

	path, err := WriteNotificationHookScript(destDir)
	if err != nil {
		t.Fatalf("WriteNotificationHookScript() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	wantMode := os.FileMode(0755)
	if info.Mode().Perm() != wantMode {
		t.Errorf("WriteNotificationHookScript() file mode = %v, want %v", info.Mode().Perm(), wantMode)
	}
}

func TestWriteNotificationHookScript_CreatesParentDirectories(t *testing.T) {
	t.Parallel()

	baseDir := t.TempDir()
	destDir := filepath.Join(baseDir, "nested", "dirs")

	path, err := WriteNotificationHookScript(destDir)
	if err != nil {
		t.Fatalf("WriteNotificationHookScript() error = %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("WriteNotificationHookScript() file does not exist at %v", path)
	}
}

func TestWritePlist_CreatesFileWithCorrectContent(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()
	binaryPath := "/usr/local/bin/claude-notifier"

	path, err := WritePlist(destDir, binaryPath)
	if err != nil {
		t.Fatalf("WritePlist() error = %v", err)
	}

	wantPath := filepath.Join(destDir, "com.dee.claude-notifier.plist")
	if path != wantPath {
		t.Errorf("WritePlist() path = %v, want %v", path, wantPath)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	if !strings.Contains(string(content), binaryPath) {
		t.Error("WritePlist() content should contain binary path")
	}
}

func TestWritePlist_SubstitutesBinaryPathPlaceholder(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()
	binaryPath := "/custom/path/to/binary"

	path, err := WritePlist(destDir, binaryPath)
	if err != nil {
		t.Fatalf("WritePlist() error = %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	contentStr := string(content)

	if strings.Contains(contentStr, "__BINARY_PATH__") {
		t.Error("WritePlist() content should not contain __BINARY_PATH__ placeholder")
	}

	if !strings.Contains(contentStr, binaryPath) {
		t.Errorf("WritePlist() content should contain substituted path %v", binaryPath)
	}
}

func TestWritePlist_CreatesParentDirectories(t *testing.T) {
	t.Parallel()

	baseDir := t.TempDir()
	destDir := filepath.Join(baseDir, "nested", "dirs")

	path, err := WritePlist(destDir, "/usr/local/bin/claude-notifier")
	if err != nil {
		t.Fatalf("WritePlist() error = %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("WritePlist() file does not exist at %v", path)
	}
}

func TestWritePlist_ReturnsErrorOnMkdirAllFailure(t *testing.T) {
	t.Parallel()

	// Use a non-existent parent with a file blocker
	baseDir := t.TempDir()
	blocker := filepath.Join(baseDir, "blocker")
	if err := os.WriteFile(blocker, []byte{}, 0644); err != nil {
		t.Fatalf("failed to create blocker file: %v", err)
	}
	// Try to write into a path where a file exists (can't create dir)
	destDir := filepath.Join(blocker, "subdir")

	_, err := WritePlist(destDir, "/usr/local/bin/claude-notifier")
	if err == nil {
		t.Error("WritePlist() expected error for invalid path, got nil")
	}
}

func TestWritePlist_ReturnsErrorOnWriteFileFailure(t *testing.T) {
	t.Parallel()

	// Create a directory where the file should be written.
	// WriteFile will fail because it cannot write to a directory.
	destDir := t.TempDir()
	targetPath := filepath.Join(destDir, "com.dee.claude-notifier.plist")
	if err := os.Mkdir(targetPath, 0755); err != nil {
		t.Fatalf("failed to create blocking directory: %v", err)
	}

	_, err := WritePlist(destDir, "/usr/local/bin/claude-notifier")
	if err == nil {
		t.Error("WritePlist() expected error when target is a directory, got nil")
	}
}

func TestWritePlist_SetsCorrectFilePermissions(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()

	path, err := WritePlist(destDir, "/usr/local/bin/claude-notifier")
	if err != nil {
		t.Fatalf("WritePlist() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	wantMode := os.FileMode(0644)
	if info.Mode().Perm() != wantMode {
		t.Errorf("WritePlist() file mode = %v, want %v", info.Mode().Perm(), wantMode)
	}
}

func TestWritePlist_ProducesValidXMLAfterSubstitution(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()
	binaryPath := "/usr/local/bin/claude-notifier"

	path, err := WritePlist(destDir, binaryPath)
	if err != nil {
		t.Fatalf("WritePlist() error = %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	contentStr := string(content)

	// Verify XML structure is intact after substitution
	requiredElements := []string{
		"<?xml version",
		"<plist",
		"</plist>",
		"<dict>",
		"</dict>",
		"<key>Label</key>",
		"<key>ProgramArguments</key>",
		"<array>",
		"</array>",
	}

	for _, elem := range requiredElements {
		if !strings.Contains(contentStr, elem) {
			t.Errorf("WritePlist() output missing required XML element: %q", elem)
		}
	}

	// Verify binary path was properly substituted in correct location
	expectedSubstitution := "<string>" + binaryPath + "</string>"
	if !strings.Contains(contentStr, expectedSubstitution) {
		t.Error("WritePlist() binary path not properly substituted within <string> tags")
	}
}
