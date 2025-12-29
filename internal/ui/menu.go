package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/jeroendee/claude-notifier/internal/notification"
)

// Menu manages the system tray menu showing notification history.
type Menu struct {
	systray   *Systray
	store     *notification.Store
	builtMenu *fyne.Menu
}

// NewMenu creates a new Menu with the given systray and notification store.
func NewMenu(systray *Systray, store *notification.Store) *Menu {
	return &Menu{
		systray: systray,
		store:   store,
	}
}

// Build constructs the system tray menu from the notification store.
func (m *Menu) Build() *fyne.Menu {
	var items []*fyne.MenuItem

	// Header with unread count
	unreadCount := m.store.UnreadCount()
	header := fyne.NewMenuItem(fmt.Sprintf("Notifications (%d unread)", unreadCount), nil)
	header.Disabled = true
	items = append(items, header)

	// Separator after header
	items = append(items, fyne.NewMenuItemSeparator())

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
		item := fyne.NewMenuItem(label, func() {
			m.store.MarkRead(notificationID)
		})
		items = append(items, item)
	}

	// Separator before actions
	items = append(items, fyne.NewMenuItemSeparator())

	// Action items
	markAllRead := fyne.NewMenuItem("Mark All Read", func() {
		m.store.MarkAllRead()
	})
	items = append(items, markAllRead)

	clearHistory := fyne.NewMenuItem("Clear History", func() {
		m.store.Clear()
	})
	items = append(items, clearHistory)

	// Separator before quit
	items = append(items, fyne.NewMenuItemSeparator())

	quit := fyne.NewMenuItem("Quit", func() {
		if m.systray != nil && m.systray.App() != nil {
			m.systray.App().Quit()
		}
	})
	items = append(items, quit)

	m.builtMenu = fyne.NewMenu("", items...)
	return m.builtMenu
}

// Refresh rebuilds the menu from the current store state.
func (m *Menu) Refresh() {
	m.Build()
}
