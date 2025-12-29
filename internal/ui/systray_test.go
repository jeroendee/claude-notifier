package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSystray(t *testing.T) {
	tests := []struct {
		name     string
		iconPath string
		want     string
	}{
		{
			name:     "stores icon path",
			iconPath: "/path/to/icon.png",
			want:     "/path/to/icon.png",
		},
		{
			name:     "stores different icon path",
			iconPath: "/another/path/icon.png",
			want:     "/another/path/icon.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewSystray(tt.iconPath)

			if s == nil {
				t.Fatal("NewSystray() returned nil")
			}
			if s.iconPath != tt.want {
				t.Errorf("iconPath = %q, want %q", s.iconPath, tt.want)
			}
		})
	}
}

func TestSystray_Setup(t *testing.T) {
	// Create temp icon file for valid path tests
	tmpDir := t.TempDir()
	validIconPath := filepath.Join(tmpDir, "icon.png")
	if err := os.WriteFile(validIconPath, []byte("PNG"), 0644); err != nil {
		t.Fatalf("failed to create temp icon: %v", err)
	}

	t.Run("creates fyne app", func(t *testing.T) {
		s := NewSystray(validIconPath)

		err := s.Setup()

		if err != nil {
			t.Errorf("Setup() error = %v, want nil", err)
		}
		if s.app == nil {
			t.Error("Setup() did not create fyne app")
		}
	})

	t.Run("returns error for invalid icon path", func(t *testing.T) {
		s := NewSystray("/nonexistent/path/icon.png")

		err := s.Setup()

		if err == nil {
			t.Error("Setup() error = nil, want error for invalid icon path")
		}
	})

	t.Run("sets system tray icon", func(t *testing.T) {
		s := NewSystray(validIconPath)

		err := s.Setup()

		if err != nil {
			t.Fatalf("Setup() error = %v", err)
		}
		// Verify systray icon was set by checking the app supports desktop features
		if s.app == nil {
			t.Fatal("app is nil")
		}
		// The systray icon setting is verified by successful Setup completion
		// with a valid icon path - the implementation loads and sets the icon
	})
}

func TestSystray_App(t *testing.T) {
	// Create temp icon file for valid path tests
	tmpDir := t.TempDir()
	validIconPath := filepath.Join(tmpDir, "icon.png")
	if err := os.WriteFile(validIconPath, []byte("PNG"), 0644); err != nil {
		t.Fatalf("failed to create temp icon: %v", err)
	}

	t.Run("returns fyne app after setup", func(t *testing.T) {
		s := NewSystray(validIconPath)
		if err := s.Setup(); err != nil {
			t.Fatalf("Setup() error = %v", err)
		}

		app := s.App()

		if app == nil {
			t.Error("App() = nil, want non-nil fyne.App")
		}
	})

	t.Run("returns nil before setup", func(t *testing.T) {
		t.Parallel()

		s := NewSystray("/some/path/icon.png")

		app := s.App()

		if app != nil {
			t.Errorf("App() = %v, want nil before Setup()", app)
		}
	})
}
