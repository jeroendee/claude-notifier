package ui

import (
	"sync"
	"testing"

	"github.com/jeroendee/claude-notifier/internal/notification"
)

// mockPlayer implements Player interface for testing.
type mockPlayer struct {
	mu     sync.RWMutex
	muted  bool
	setCnt int
}

func (m *mockPlayer) IsMuted() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.muted
}

func (m *mockPlayer) SetMuted(muted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.muted = muted
	m.setCnt++
}

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

func TestMenuHasPlayerField(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)

	// Verify player field exists (will be nil before SetPlayer is called)
	if menu.player != nil {
		t.Error("player should be nil initially")
	}
}

func TestMenuHasMuteItemField(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)

	// Verify muteItem field exists (will be nil before Build is called)
	if menu.muteItem != nil {
		t.Error("muteItem should be nil before Build is called")
	}
}

func TestMenuSetPlayer(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)
	player := &mockPlayer{}

	menu.SetPlayer(player)

	if menu.player != player {
		t.Error("SetPlayer did not set player")
	}
}

func TestMenuSetPlayerNil(t *testing.T) {
	t.Parallel()

	systray := &Systray{}
	store := notification.NewStore()
	menu := NewMenu(systray, store)
	player := &mockPlayer{}

	menu.SetPlayer(player)
	menu.SetPlayer(nil)

	if menu.player != nil {
		t.Error("SetPlayer(nil) did not clear player")
	}
}

func TestMenuToggleMuteFromUnmuted(t *testing.T) {
	t.Parallel()

	menu := &Menu{}
	player := &mockPlayer{muted: false}
	menu.player = player

	menu.toggleMute()

	if !player.muted {
		t.Error("toggleMute() did not mute player")
	}
	if player.setCnt != 1 {
		t.Errorf("SetMuted called %d times, want 1", player.setCnt)
	}
}

func TestMenuToggleMuteFromMuted(t *testing.T) {
	t.Parallel()

	menu := &Menu{}
	player := &mockPlayer{muted: true}
	menu.player = player

	menu.toggleMute()

	if player.muted {
		t.Error("toggleMute() did not unmute player")
	}
	if player.setCnt != 1 {
		t.Errorf("SetMuted called %d times, want 1", player.setCnt)
	}
}

func TestMenuToggleMuteNilPlayer(t *testing.T) {
	t.Parallel()

	menu := &Menu{}
	menu.player = nil

	// Should not panic
	menu.toggleMute()
}

func TestMenuMuteLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		muted    bool
		wantText string
	}{
		{"unmuted shows Mute Sound", false, "Mute Sound"},
		{"muted shows Unmute Sound", true, "Unmute Sound"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			player := &mockPlayer{muted: tt.muted}
			got := muteLabel(player)
			if got != tt.wantText {
				t.Errorf("muteLabel() = %q, want %q", got, tt.wantText)
			}
		})
	}
}

func TestMenuMuteLabelNilPlayer(t *testing.T) {
	t.Parallel()

	got := muteLabel(nil)
	if got != "Mute Sound" {
		t.Errorf("muteLabel(nil) = %q, want %q", got, "Mute Sound")
	}
}

func TestMenuUpdateMuteLabelNilMuteItem(t *testing.T) {
	t.Parallel()

	menu := &Menu{}
	player := &mockPlayer{muted: false}
	menu.player = player
	menu.muteItem = nil

	// Should not panic
	menu.updateMuteLabel()
}
