package ui

import (
	"github.com/getlantern/systray"

	"github.com/jeroendee/claude-notifier/internal/assets"
)

// Systray manages the system tray icon and menu bar presence.
type Systray struct {
	iconData        []byte
	alertIconData   []byte
	currentIconData []byte
	menu            *Menu
	quitCh          chan struct{}
}

// NewSystray creates a new Systray with embedded icon assets.
func NewSystray() *Systray {
	return &Systray{
		iconData:      assets.Icon,
		alertIconData: assets.IconAlert,
		quitCh:        make(chan struct{}),
	}
}

// Setup initializes the systray. Icons are already embedded.
func (s *Systray) Setup() error {
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
