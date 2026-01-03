package scripts

import (
	"bytes"
	"strings"
	"testing"
)

func TestHookScript(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		check func(t *testing.T, content []byte)
	}{
		{
			name: "returns non-empty bytes",
			check: func(t *testing.T, content []byte) {
				if len(content) == 0 {
					t.Error("HookScript should return non-empty []byte")
				}
			},
		},
		{
			name: "starts with correct shebang",
			check: func(t *testing.T, content []byte) {
				wantShebang := []byte("#!/bin/bash\n")
				if !bytes.HasPrefix(content, wantShebang) {
					t.Errorf("HookScript should start with shebang, got prefix: %q", content[:min(len(content), 20)])
				}
			},
		},
		{
			name: "contains curl command",
			check: func(t *testing.T, content []byte) {
				if !bytes.Contains(content, []byte("curl")) {
					t.Error("HookScript should contain curl command")
				}
			},
		},
		{
			name: "uses jq for JSON construction to prevent injection",
			check: func(t *testing.T, content []byte) {
				if !bytes.Contains(content, []byte("jq -n")) {
					t.Error("HookScript should use 'jq -n' for safe JSON construction")
				}
			},
		},
		{
			name: "does not use string interpolation for JSON",
			check: func(t *testing.T, content []byte) {
				// The vulnerable pattern is '{"message": "'"$MESSAGE"'"}'
				if bytes.Contains(content, []byte(`'{"message"`)) {
					t.Error("HookScript should not use string interpolation for JSON (vulnerable to injection)")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.check(t, HookScript)
		})
	}
}

func TestNotificationHookScript(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		check func(t *testing.T, content []byte)
	}{
		{
			name: "returns non-empty bytes",
			check: func(t *testing.T, content []byte) {
				if len(content) == 0 {
					t.Error("NotificationHookScript should return non-empty []byte")
				}
			},
		},
		{
			name: "starts with correct shebang",
			check: func(t *testing.T, content []byte) {
				wantShebang := []byte("#!/bin/bash\n")
				if !bytes.HasPrefix(content, wantShebang) {
					t.Errorf("NotificationHookScript should start with shebang, got prefix: %q", content[:min(len(content), 20)])
				}
			},
		},
		{
			name: "contains curl command",
			check: func(t *testing.T, content []byte) {
				if !bytes.Contains(content, []byte("curl")) {
					t.Error("NotificationHookScript should contain curl command")
				}
			},
		},
		{
			name: "uses jq for JSON construction to prevent injection",
			check: func(t *testing.T, content []byte) {
				if !bytes.Contains(content, []byte("jq -n")) {
					t.Error("NotificationHookScript should use 'jq -n' for safe JSON construction")
				}
			},
		},
		{
			name: "does not use string interpolation for JSON",
			check: func(t *testing.T, content []byte) {
				// The vulnerable pattern is '{"message": "'"$MESSAGE"'"}'
				if bytes.Contains(content, []byte(`'{"message"`)) {
					t.Error("NotificationHookScript should not use string interpolation for JSON (vulnerable to injection)")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.check(t, NotificationHookScript)
		})
	}
}

func TestPlistTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		check func(t *testing.T, content []byte)
	}{
		{
			name: "returns non-empty bytes",
			check: func(t *testing.T, content []byte) {
				if len(content) == 0 {
					t.Error("PlistTemplate should return non-empty []byte")
				}
			},
		},
		{
			name: "contains XML declaration",
			check: func(t *testing.T, content []byte) {
				if !bytes.HasPrefix(content, []byte("<?xml")) {
					t.Error("PlistTemplate should start with XML declaration")
				}
			},
		},
		{
			name: "contains binary path placeholder",
			check: func(t *testing.T, content []byte) {
				if !bytes.Contains(content, []byte("__BINARY_PATH__")) {
					t.Error("PlistTemplate should contain __BINARY_PATH__ placeholder")
				}
			},
		},
		{
			name: "has valid plist structure",
			check: func(t *testing.T, content []byte) {
				str := string(content)
				required := []string{
					"<plist",
					"</plist>",
					"<dict>",
					"</dict>",
					"<key>Label</key>",
					"<key>ProgramArguments</key>",
				}
				for _, req := range required {
					if !strings.Contains(str, req) {
						t.Errorf("PlistTemplate should contain %q", req)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.check(t, PlistTemplate)
		})
	}
}
