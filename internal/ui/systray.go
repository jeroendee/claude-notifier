package ui

import (
	"fmt"
	"os"

	"fyne.io/systray"
)

// Systray manages the system tray icon and menu bar presence.
type Systray struct {
	iconPath string
	iconData []byte
	menu     *Menu
	quitCh   chan struct{}
}

// NewSystray creates a new Systray with the specified icon path.
func NewSystray(iconPath string) *Systray {
	return &Systray{
		iconPath: iconPath,
		quitCh:   make(chan struct{}),
	}
}

// Setup validates the icon file exists and loads it.
func (s *Systray) Setup() error {
	iconData, err := os.ReadFile(s.iconPath)
	if err != nil {
		return fmt.Errorf("icon file not found: %w", err)
	}
	s.iconData = iconData
	return nil
}

// SetMenu sets the menu to be displayed in the systray.
func (s *Systray) SetMenu(menu *Menu) {
	s.menu = menu
}

// Run starts the systray event loop (blocking).
func (s *Systray) Run() {
	systray.Run(s.onReady, s.onExit)
}

// Quit signals the systray to quit.
func (s *Systray) Quit() {
	systray.Quit()
}

func (s *Systray) onReady() {
	systray.SetIcon(s.iconData)
	systray.SetTooltip("Claude Notifier")

	if s.menu != nil {
		s.menu.Build()
	}
}

func (s *Systray) onExit() {
	close(s.quitCh)
}
