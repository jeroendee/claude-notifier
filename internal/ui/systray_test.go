package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSystray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		iconPath      string
		alertIconPath string
	}{
		{
			name:          "stores both icon paths",
			iconPath:      "/path/to/icon.png",
			alertIconPath: "/path/to/icon-alert.png",
		},
		{
			name:          "stores different icon paths",
			iconPath:      "/another/path/icon.ico",
			alertIconPath: "/another/path/icon-alert.ico",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewSystray(tt.iconPath, tt.alertIconPath)

			if s == nil {
				t.Fatal("NewSystray returned nil")
			}
			if s.iconPath != tt.iconPath {
				t.Errorf("iconPath = %q, want %q", s.iconPath, tt.iconPath)
			}
			if s.alertIconPath != tt.alertIconPath {
				t.Errorf("alertIconPath = %q, want %q", s.alertIconPath, tt.alertIconPath)
			}
			if s.quitCh == nil {
				t.Error("quitCh should be initialized")
			}
		})
	}
}

func TestSystray_Setup(t *testing.T) {
	t.Parallel()

	t.Run("returns error for invalid main icon path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		alertIconPath := filepath.Join(tmpDir, "alert-icon.png")
		if err := os.WriteFile(alertIconPath, []byte{0x89}, 0644); err != nil {
			t.Fatalf("failed to create alert icon: %v", err)
		}

		s := NewSystray("/nonexistent/icon.png", alertIconPath)
		err := s.Setup()

		if err == nil {
			t.Error("Setup should return error for invalid main icon path")
		}
	})

	t.Run("returns error for invalid alert icon path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		iconPath := filepath.Join(tmpDir, "icon.png")
		if err := os.WriteFile(iconPath, []byte{0x89}, 0644); err != nil {
			t.Fatalf("failed to create main icon: %v", err)
		}

		s := NewSystray(iconPath, "/nonexistent/alert-icon.png")
		err := s.Setup()

		if err == nil {
			t.Error("Setup should return error for invalid alert icon path")
		}
	})

	t.Run("loads both icon data for valid paths", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		iconPath := filepath.Join(tmpDir, "test-icon.png")
		alertIconPath := filepath.Join(tmpDir, "test-alert-icon.png")
		iconData := []byte{0x89, 0x50, 0x4E, 0x47}
		alertIconData := []byte{0xFF, 0xD8, 0xFF, 0xE0}

		if err := os.WriteFile(iconPath, iconData, 0644); err != nil {
			t.Fatalf("failed to create test icon: %v", err)
		}
		if err := os.WriteFile(alertIconPath, alertIconData, 0644); err != nil {
			t.Fatalf("failed to create alert icon: %v", err)
		}

		s := NewSystray(iconPath, alertIconPath)
		err := s.Setup()

		if err != nil {
			t.Errorf("Setup returned error: %v", err)
		}
		if len(s.iconData) != len(iconData) {
			t.Errorf("iconData length = %d, want %d", len(s.iconData), len(iconData))
		}
		if len(s.alertIconData) != len(alertIconData) {
			t.Errorf("alertIconData length = %d, want %d", len(s.alertIconData), len(alertIconData))
		}
	})
}

func TestSystray_SetMenu(t *testing.T) {
	t.Parallel()

	s := NewSystray("/path/to/icon.png", "/path/to/alert.png")
	menu := &Menu{}

	s.SetMenu(menu)

	if s.menu != menu {
		t.Error("SetMenu did not set menu")
	}
}

func TestSystray_SetAlertState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		hasUnread bool
		wantAlert bool
	}{
		{
			name:      "true shows alert icon",
			hasUnread: true,
			wantAlert: true,
		},
		{
			name:      "false shows normal icon",
			hasUnread: false,
			wantAlert: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			iconPath := filepath.Join(tmpDir, "icon.png")
			alertIconPath := filepath.Join(tmpDir, "alert-icon.png")
			iconData := []byte{0x89, 0x50, 0x4E, 0x47}       // normal icon
			alertIconData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // alert icon

			if err := os.WriteFile(iconPath, iconData, 0644); err != nil {
				t.Fatalf("failed to create icon: %v", err)
			}
			if err := os.WriteFile(alertIconPath, alertIconData, 0644); err != nil {
				t.Fatalf("failed to create alert icon: %v", err)
			}

			s := NewSystray(iconPath, alertIconPath)
			if err := s.Setup(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			s.SetAlertState(tt.hasUnread)

			if tt.wantAlert {
				if string(s.currentIconData) != string(alertIconData) {
					t.Error("SetAlertState(true) should set alert icon as current")
				}
			} else {
				if string(s.currentIconData) != string(iconData) {
					t.Error("SetAlertState(false) should set normal icon as current")
				}
			}
		})
	}
}
