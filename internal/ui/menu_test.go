package ui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/jeroendee/claude-notifier/internal/notification"
)

func TestNewMenu_CreatesMenuWithSystrayAndStore(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()

	menu := NewMenu(systray, store)

	if menu == nil {
		t.Fatal("NewMenu() returned nil")
	}
	if menu.systray != systray {
		t.Error("NewMenu() systray reference is incorrect")
	}
	if menu.store != store {
		t.Error("NewMenu() store reference is incorrect")
	}
}

func TestBuild_ReturnsMenuWithCorrectStructure(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	if fyneMenu == nil {
		t.Fatal("Build() returned nil")
	}

	// Find expected menu items by label
	var foundHeader, foundMarkAllRead, foundClearHistory, foundQuit bool
	for _, item := range fyneMenu.Items {
		if item.Label == "" && item.IsSeparator {
			continue // separator
		}
		if strings.HasPrefix(item.Label, "Notifications") {
			foundHeader = true
		}
		if item.Label == "Mark All Read" {
			foundMarkAllRead = true
		}
		if item.Label == "Clear History" {
			foundClearHistory = true
		}
		if item.Label == "Quit" {
			foundQuit = true
		}
	}

	if !foundHeader {
		t.Error("Build() menu missing header item")
	}
	if !foundMarkAllRead {
		t.Error("Build() menu missing 'Mark All Read' item")
	}
	if !foundClearHistory {
		t.Error("Build() menu missing 'Clear History' item")
	}
	if !foundQuit {
		t.Error("Build() menu missing 'Quit' item")
	}
}

func TestBuild_HeaderShowsZeroUnreadCount(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Find header item
	var header *fyne.MenuItem
	for _, item := range fyneMenu.Items {
		if strings.HasPrefix(item.Label, "Notifications") {
			header = item
			break
		}
	}

	if header == nil {
		t.Fatal("Build() missing header item")
	}
	if header.Label != "Notifications (0 unread)" {
		t.Errorf("Build() header = %q, want %q", header.Label, "Notifications (0 unread)")
	}
}

func TestBuild_HeaderShowsUnreadCount(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	store.Add("message 1")
	store.Add("message 2")
	store.Add("message 3")
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Find header item
	var header *fyne.MenuItem
	for _, item := range fyneMenu.Items {
		if strings.HasPrefix(item.Label, "Notifications") {
			header = item
			break
		}
	}

	if header == nil {
		t.Fatal("Build() missing header item")
	}
	if header.Label != "Notifications (3 unread)" {
		t.Errorf("Build() header = %q, want %q", header.Label, "Notifications (3 unread)")
	}
}

func TestBuild_ShowsUnreadIndicator(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	store.Add("unread message")
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Find notification item (starts with ●)
	var found bool
	for _, item := range fyneMenu.Items {
		if strings.HasPrefix(item.Label, "●") {
			found = true
			if !strings.Contains(item.Label, "unread message") {
				t.Errorf("Build() unread item = %q, want to contain 'unread message'", item.Label)
			}
			break
		}
	}

	if !found {
		t.Error("Build() missing unread indicator (●) for unread notification")
	}
}

func TestBuild_ShowsReadIndicator(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	store.Add("read message")
	notifications := store.List()
	store.MarkRead(notifications[0].ID)
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Find notification item (starts with ○)
	var found bool
	for _, item := range fyneMenu.Items {
		if strings.HasPrefix(item.Label, "○") {
			found = true
			if !strings.Contains(item.Label, "read message") {
				t.Errorf("Build() read item = %q, want to contain 'read message'", item.Label)
			}
			break
		}
	}

	if !found {
		t.Error("Build() missing read indicator (○) for read notification")
	}
}

func TestBuild_NotificationsInReverseOrder(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	store.Add("first message")
	store.Add("second message")
	store.Add("third message")
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Collect notification items (those starting with ● or ○)
	var notificationLabels []string
	for _, item := range fyneMenu.Items {
		if strings.HasPrefix(item.Label, "●") || strings.HasPrefix(item.Label, "○") {
			notificationLabels = append(notificationLabels, item.Label)
		}
	}

	if len(notificationLabels) != 3 {
		t.Fatalf("Build() notification count = %d, want 3", len(notificationLabels))
	}

	// Newest first (third message), oldest last (first message)
	if !strings.Contains(notificationLabels[0], "third message") {
		t.Errorf("Build() first notification = %q, want to contain 'third message'", notificationLabels[0])
	}
	if !strings.Contains(notificationLabels[2], "first message") {
		t.Errorf("Build() last notification = %q, want to contain 'first message'", notificationLabels[2])
	}
}

func TestBuild_EmptyStoreShowsNoNotifications(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Count notification items (those starting with ● or ○)
	notificationCount := 0
	for _, item := range fyneMenu.Items {
		if strings.HasPrefix(item.Label, "●") || strings.HasPrefix(item.Label, "○") {
			notificationCount++
		}
	}

	if notificationCount != 0 {
		t.Errorf("Build() notification count = %d, want 0 for empty store", notificationCount)
	}
}

func TestBuild_MarkAllReadActionCallsStore(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	store.Add("message 1")
	store.Add("message 2")
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Find and trigger Mark All Read action
	var markAllRead *fyne.MenuItem
	for _, item := range fyneMenu.Items {
		if item.Label == "Mark All Read" {
			markAllRead = item
			break
		}
	}

	if markAllRead == nil {
		t.Fatal("Build() missing 'Mark All Read' item")
	}

	// Trigger the action
	markAllRead.Action()

	// Verify store was updated
	if store.UnreadCount() != 0 {
		t.Errorf("Mark All Read action did not mark all as read, UnreadCount() = %d, want 0", store.UnreadCount())
	}
}

func TestBuild_ClearHistoryActionCallsStore(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	store.Add("message 1")
	store.Add("message 2")
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Find and trigger Clear History action
	var clearHistory *fyne.MenuItem
	for _, item := range fyneMenu.Items {
		if item.Label == "Clear History" {
			clearHistory = item
			break
		}
	}

	if clearHistory == nil {
		t.Fatal("Build() missing 'Clear History' item")
	}

	// Trigger the action
	clearHistory.Action()

	// Verify store was cleared
	if len(store.List()) != 0 {
		t.Errorf("Clear History action did not clear store, len = %d, want 0", len(store.List()))
	}
}

func TestBuild_ClickingNotificationMarksAsRead(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	store.Add("message to mark read")
	menu := NewMenu(systray, store)

	fyneMenu := menu.Build()

	// Find the notification item
	var notificationItem *fyne.MenuItem
	for _, item := range fyneMenu.Items {
		if strings.HasPrefix(item.Label, "●") && strings.Contains(item.Label, "message to mark read") {
			notificationItem = item
			break
		}
	}

	if notificationItem == nil {
		t.Fatal("Build() missing notification item")
	}

	// Trigger the action (click)
	notificationItem.Action()

	// Verify notification was marked as read
	notifications := store.List()
	if !notifications[0].Read {
		t.Error("Clicking notification did not mark it as read")
	}
}

func TestRefresh_RebuildsMen(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)

	// Build initial menu
	_ = menu.Build()

	// Add a notification
	store.Add("new message")

	// Refresh menu
	menu.Refresh()

	// Build again to verify refresh worked
	fyneMenu := menu.Build()

	// Should now have one notification
	var found bool
	for _, item := range fyneMenu.Items {
		if strings.Contains(item.Label, "new message") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Refresh() did not rebuild menu with new notification")
	}
}
