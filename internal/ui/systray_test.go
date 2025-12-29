package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSystray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		iconPath string
	}{
		{
			name:     "stores icon path",
			iconPath: "/path/to/icon.png",
		},
		{
			name:     "stores different icon path",
			iconPath: "/another/path/icon.ico",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewSystray(tt.iconPath)

			if s == nil {
				t.Fatal("NewSystray returned nil")
			}
			if s.iconPath != tt.iconPath {
				t.Errorf("iconPath = %q, want %q", s.iconPath, tt.iconPath)
			}
			if s.quitCh == nil {
				t.Error("quitCh should be initialized")
			}
		})
	}
}

func TestSystray_Setup(t *testing.T) {
	t.Parallel()

	t.Run("returns error for invalid icon path", func(t *testing.T) {
		t.Parallel()

		s := NewSystray("/nonexistent/icon.png")
		err := s.Setup()

		if err == nil {
			t.Error("Setup should return error for invalid path")
		}
	})

	t.Run("loads icon data for valid path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		iconPath := filepath.Join(tmpDir, "test-icon.png")
		iconData := []byte{0x89, 0x50, 0x4E, 0x47}
		if err := os.WriteFile(iconPath, iconData, 0644); err != nil {
			t.Fatalf("failed to create test icon: %v", err)
		}

		s := NewSystray(iconPath)
		err := s.Setup()

		if err != nil {
			t.Errorf("Setup returned error: %v", err)
		}
		if len(s.iconData) != len(iconData) {
			t.Errorf("iconData length = %d, want %d", len(s.iconData), len(iconData))
		}
	})
}

func TestSystray_SetMenu(t *testing.T) {
	t.Parallel()

	s := NewSystray("/path/to/icon.png")
	menu := &Menu{}

	s.SetMenu(menu)

	if s.menu != menu {
		t.Error("SetMenu did not set menu")
	}
}
