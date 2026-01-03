package ui

import (
	"testing"

	"github.com/jeroendee/claude-notifier/internal/notification"
)

func TestNewMenu(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()

	menu := NewMenu(systray, store)

	if menu == nil {
		t.Fatal("NewMenu returned nil")
	}
	if menu.systray != systray {
		t.Error("NewMenu did not set systray")
	}
	if menu.store != store {
		t.Error("NewMenu did not set store")
	}
}

func TestMenuHasAboutField(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)

	// Verify about field exists (will be nil before Build is called)
	// This test ensures the Menu struct has the about field defined
	if menu.about != nil {
		t.Error("about should be nil before Build is called")
	}
}

// Note: Build() and Refresh() methods use fyne.io/systray global state
// and cannot be unit tested in isolation. These are tested through
// integration testing by running the actual application.

func TestNotificationItemConcurrentAccess(t *testing.T) {
	t.Parallel()

	ni := &notificationItem{}

	// Simulate concurrent read/write pattern from Build() and Refresh()
	done := make(chan struct{})
	const iterations = 1000

	// Writer goroutine (simulates Refresh)
	go func() {
		for i := 0; i < iterations; i++ {
			ni.SetID("test-id")
			ni.SetID("")
		}
		close(done)
	}()

	// Reader goroutine (simulates click handler in Build)
	for i := 0; i < iterations; i++ {
		_ = ni.GetID()
	}

	<-done
}

func TestNotificationItemGetID(t *testing.T) {
	t.Parallel()

	ni := &notificationItem{}

	if got := ni.GetID(); got != "" {
		t.Errorf("GetID() = %q, want empty string", got)
	}

	ni.SetID("abc123")

	if got := ni.GetID(); got != "abc123" {
		t.Errorf("GetID() = %q, want %q", got, "abc123")
	}
}

func TestNotificationItemSetID(t *testing.T) {
	t.Parallel()

	ni := &notificationItem{}

	ni.SetID("first")
	if got := ni.GetID(); got != "first" {
		t.Errorf("after SetID(first), GetID() = %q, want %q", got, "first")
	}

	ni.SetID("second")
	if got := ni.GetID(); got != "second" {
		t.Errorf("after SetID(second), GetID() = %q, want %q", got, "second")
	}

	ni.SetID("")
	if got := ni.GetID(); got != "" {
		t.Errorf("after SetID(\"\"), GetID() = %q, want empty", got)
	}
}
