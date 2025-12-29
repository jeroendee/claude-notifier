package ui

import (
	"fmt"

	"fyne.io/systray"
	"github.com/jeroendee/claude-notifier/internal/notification"
)

// Menu manages the system tray menu showing notification history.
type Menu struct {
	systray *Systray
	store   *notification.Store
}

// NewMenu creates a new Menu with the given systray and notification store.
func NewMenu(systray *Systray, store *notification.Store) *Menu {
	return &Menu{
		systray: systray,
		store:   store,
	}
}

// Build constructs the system tray menu from the notification store.
func (m *Menu) Build() {
	systray.ResetMenu()

	// Header with unread count
	unreadCount := m.store.UnreadCount()
	header := systray.AddMenuItem(fmt.Sprintf("Notifications (%d unread)", unreadCount), "")
	header.Disable()

	systray.AddSeparator()

	// Notification history items (newest first)
	notifications := m.store.List()
	for i := len(notifications) - 1; i >= 0; i-- {
		n := notifications[i]
		indicator := "●" // unread
		if n.Read {
			indicator = "○" // read
		}
		label := fmt.Sprintf("%s %s", indicator, n.Message)
		notificationID := n.ID
		item := systray.AddMenuItem(label, "Click to mark as read")

		go func(id string, menuItem *systray.MenuItem) {
			for range menuItem.ClickedCh {
				m.store.MarkRead(id)
			}
		}(notificationID, item)
	}

	systray.AddSeparator()

	// Action items
	markAllRead := systray.AddMenuItem("Mark All Read", "Mark all notifications as read")
	go func() {
		for range markAllRead.ClickedCh {
			m.store.MarkAllRead()
		}
	}()

	clearHistory := systray.AddMenuItem("Clear History", "Remove all notifications")
	go func() {
		for range clearHistory.ClickedCh {
			m.store.Clear()
		}
	}()

	systray.AddSeparator()

	quit := systray.AddMenuItem("Quit", "Quit the application")
	go func() {
		for range quit.ClickedCh {
			if m.systray != nil {
				m.systray.Quit()
			}
		}
	}()
}

// Refresh rebuilds the menu from the current store state.
func (m *Menu) Refresh() {
	m.Build()
}
