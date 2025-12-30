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
