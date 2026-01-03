// Package ui provides the system tray user interface components.
//
// # Goroutine Lifecycle
//
// Menu spawns goroutines in Build() that range over systray MenuItem.ClickedCh
// channels to handle click events. These goroutines are designed to run for the
// lifetime of the application.
//
// The systray package (github.com/getlantern/systray) does not close ClickedCh
// channels on Quit(). Instead, goroutine cleanup relies on process termination:
// when systray.Quit() is called, the systray event loop exits, Run() returns,
// and the process terminates, killing all goroutines.
//
// This is intentional behavior - adding explicit done channel cleanup would add
// complexity for no practical benefit, since the process terminates immediately
// after systray shutdown.
package ui

import (
	"fmt"
	"sync"

	"github.com/getlantern/systray"
	"github.com/jeroendee/claude-notifier/internal/notification"
	"github.com/jeroendee/claude-notifier/internal/version"
)

const maxNotificationItems = 20

// Player abstracts sound playback control for dependency injection.
type Player interface {
	IsMuted() bool
	SetMuted(bool)
}

// notificationItem holds a menu item and its associated notification ID.
// Access to id is protected by mu for concurrent access from click handlers
// and Refresh().
type notificationItem struct {
	menuItem *systray.MenuItem
	mu       sync.RWMutex
	id       string
}

// GetID returns the notification ID safely for concurrent access.
func (ni *notificationItem) GetID() string {
	ni.mu.RLock()
	defer ni.mu.RUnlock()
	return ni.id
}

// SetID updates the notification ID safely for concurrent access.
func (ni *notificationItem) SetID(id string) {
	ni.mu.Lock()
	defer ni.mu.Unlock()
	ni.id = id
}

// Menu manages the system tray menu showing notification history.
type Menu struct {
	systray           *Systray
	store             *notification.Store
	player            Player
	header            *systray.MenuItem
	notificationItems []*notificationItem
	muteItem          *systray.MenuItem
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

// SetPlayer sets the player for mute control.
func (m *Menu) SetPlayer(p Player) {
	m.player = p
}

// toggleMute flips the player's muted state.
func (m *Menu) toggleMute() {
	if m.player == nil {
		return
	}
	m.player.SetMuted(!m.player.IsMuted())
}

// muteLabel returns the appropriate label based on player mute state.
func muteLabel(p Player) string {
	if p != nil && p.IsMuted() {
		return "Unmute Sound"
	}
	return "Mute Sound"
}

// updateMuteLabel updates the mute menu item title based on player state.
func (m *Menu) updateMuteLabel() {
	if m.muteItem == nil {
		return
	}
	m.muteItem.SetTitle(muteLabel(m.player))
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

	// Pre-create notification item slots (hidden initially).
	// Each item gets a goroutine that handles clicks. These goroutines run
	// until process termination (see package doc for lifecycle details).
	for i := 0; i < maxNotificationItems; i++ {
		item := systray.AddMenuItem("", "Click to mark as read")
		item.Hide()
		ni := &notificationItem{menuItem: item, id: ""}
		m.notificationItems = append(m.notificationItems, ni)

		go func(ni *notificationItem) {
			for range ni.menuItem.ClickedCh {
				if id := ni.GetID(); id != "" {
					m.store.MarkRead(id)
				}
			}
		}(ni)
	}

	systray.AddSeparator()

	// Action items. Goroutines run until process termination.
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

	// Mute sound toggle. Goroutine runs until process termination.
	m.muteItem = systray.AddMenuItem(muteLabel(m.player), "Toggle notification sound")
	go func() {
		for range m.muteItem.ClickedCh {
			m.toggleMute()
			m.updateMuteLabel()
		}
	}()

	systray.AddSeparator()

	// About item with version (disabled, informational only)
	m.about = systray.AddMenuItem(fmt.Sprintf("Version (%s)", version.Get().Version), "")
	m.about.Disable()

	// Quit handler goroutine runs until process termination.
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
		ni.SetID(n.ID)
		ni.menuItem.SetTitle(label)
		ni.menuItem.Show()
		itemIndex++
	}

	// Hide unused notification items
	for ; itemIndex < maxNotificationItems; itemIndex++ {
		ni := m.notificationItems[itemIndex]
		ni.SetID("")
		ni.menuItem.Hide()
	}
}
