package ui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

// Systray manages the system tray icon and menu bar presence.
type Systray struct {
	app      fyne.App
	iconPath string
}

// NewSystray creates a new Systray with the specified icon path.
func NewSystray(iconPath string) *Systray {
	return &Systray{
		iconPath: iconPath,
	}
}

// Setup creates the Fyne app and sets the system tray icon.
func (s *Systray) Setup() error {
	// Verify icon file exists
	iconData, err := os.ReadFile(s.iconPath)
	if err != nil {
		return fmt.Errorf("icon file not found: %w", err)
	}

	s.app = app.New()

	// Set system tray icon if desktop app is supported
	if deskApp, ok := s.app.(desktop.App); ok {
		icon := fyne.NewStaticResource("icon", iconData)
		deskApp.SetSystemTrayIcon(icon)
	}

	return nil
}

// App returns the Fyne app for menu setup.
func (s *Systray) App() fyne.App {
	return s.app
}

// Run runs the Fyne app event loop (blocking).
func (s *Systray) Run() {
	if s.app != nil {
		s.app.Run()
	}
}
