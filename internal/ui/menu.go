package ui

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/jeroendee/claude-notifier/internal/notification"
	"github.com/jeroendee/claude-notifier/internal/version"
)

const maxNotificationItems = 20

// notificationItem holds a menu item and its associated notification ID.
type notificationItem struct {
	menuItem *systray.MenuItem
	id       string
}

// Menu manages the system tray menu showing notification history.
type Menu struct {
	systray           *Systray
	store             *notification.Store
	header            *systray.MenuItem
	notificationItems []*notificationItem
	markAllRead       *systray.MenuItem
	clearHistory      *systray.MenuItem
	about             *systray.MenuItem
	quit              *systray.MenuItem
	built             bool
}

// NewMenu creates a new Menu with the given systray and notification store.
func NewMenu(systray *Systray, store *notification.Store) *Menu {
	return &Menu{
		systray:           systray,
		store:             store,
		notificationItems: make([]*notificationItem, 0, maxNotificationItems),
	}
}

// Build constructs the system tray menu. Called once during initialization.
func (m *Menu) Build() {
	if m.built {
		return
	}

	// Header with unread count
	m.header = systray.AddMenuItem("Notifications (0 unread)", "")
	m.header.Disable()

	systray.AddSeparator()

	// Pre-create notification item slots (hidden initially)
	for i := 0; i < maxNotificationItems; i++ {
		item := systray.AddMenuItem("", "Click to mark as read")
		item.Hide()
		ni := &notificationItem{menuItem: item, id: ""}
		m.notificationItems = append(m.notificationItems, ni)

		// Click handler - uses closure to access the item's current ID
		go func(ni *notificationItem) {
			for range ni.menuItem.ClickedCh {
				if ni.id != "" {
					m.store.MarkRead(ni.id)
				}
			}
		}(ni)
	}

	systray.AddSeparator()

	// Action items
	m.markAllRead = systray.AddMenuItem("Mark All Read", "Mark all notifications as read")
	go func() {
		for range m.markAllRead.ClickedCh {
			m.store.MarkAllRead()
		}
	}()

	m.clearHistory = systray.AddMenuItem("Clear History", "Remove all notifications")
	go func() {
		for range m.clearHistory.ClickedCh {
			m.store.Clear()
		}
	}()

	systray.AddSeparator()

	// About item with version (disabled, informational only)
	m.about = systray.AddMenuItem(fmt.Sprintf("About (%s)", version.Get().Version), "")
	m.about.Disable()

	m.quit = systray.AddMenuItem("Quit", "Quit the application")
	go func() {
		for range m.quit.ClickedCh {
			if m.systray != nil {
				m.systray.Quit()
			}
		}
	}()

	m.built = true

	// Initial refresh to set correct values
	m.Refresh()
}

// Refresh updates the menu to reflect the current store state.
func (m *Menu) Refresh() {
	if !m.built {
		return
	}

	// Update header with unread count
	unreadCount := m.store.UnreadCount()
	m.header.SetTitle(fmt.Sprintf("Notifications (%d unread)", unreadCount))

	// Get notifications (newest first)
	notifications := m.store.List()

	// Update notification items
	itemIndex := 0
	for i := len(notifications) - 1; i >= 0 && itemIndex < maxNotificationItems; i-- {
		n := notifications[i]
		indicator := "●" // unread
		if n.Read {
			indicator = "○" // read
		}
		label := fmt.Sprintf("%s %s", indicator, n.Message)

		ni := m.notificationItems[itemIndex]
		ni.id = n.ID
		ni.menuItem.SetTitle(label)
		ni.menuItem.Show()
		itemIndex++
	}

	// Hide unused notification items
	for ; itemIndex < maxNotificationItems; itemIndex++ {
		ni := m.notificationItems[itemIndex]
		ni.id = ""
		ni.menuItem.Hide()
	}
}
