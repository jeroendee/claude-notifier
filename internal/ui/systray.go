package ui

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
)

// Systray manages the system tray icon and menu bar presence.
type Systray struct {
	iconPath        string
	alertIconPath   string
	iconData        []byte
	alertIconData   []byte
	currentIconData []byte
	menu            *Menu
	quitCh          chan struct{}
}

// NewSystray creates a new Systray with the specified icon paths.
func NewSystray(iconPath, alertIconPath string) *Systray {
	return &Systray{
		iconPath:      iconPath,
		alertIconPath: alertIconPath,
		quitCh:        make(chan struct{}),
	}
}

// Setup validates icon files exist and loads them.
func (s *Systray) Setup() error {
	iconData, err := os.ReadFile(s.iconPath)
	if err != nil {
		return fmt.Errorf("icon file not found: %w", err)
	}
	s.iconData = iconData

	alertIconData, err := os.ReadFile(s.alertIconPath)
	if err != nil {
		return fmt.Errorf("alert icon file not found: %w", err)
	}
	s.alertIconData = alertIconData

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

// SetAlertState switches between normal and alert icons based on unread state.
func (s *Systray) SetAlertState(hasUnread bool) {
	if hasUnread {
		s.currentIconData = s.alertIconData
	} else {
		s.currentIconData = s.iconData
	}
	systray.SetIcon(s.currentIconData)
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
