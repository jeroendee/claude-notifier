package ui

import (
	"testing"

	"github.com/jeroendee/claude-notifier/internal/assets"
)

func TestNewSystray(t *testing.T) {
	t.Parallel()

	s := NewSystray()

	if s == nil {
		t.Fatal("NewSystray returned nil")
	}
	if len(s.iconData) == 0 {
		t.Error("iconData should be initialized with embedded icon")
	}
	if len(s.alertIconData) == 0 {
		t.Error("alertIconData should be initialized with embedded alert icon")
	}
	if s.quitCh == nil {
		t.Error("quitCh should be initialized")
	}
}

func TestNewSystray_UsesEmbeddedAssets(t *testing.T) {
	t.Parallel()

	s := NewSystray()

	if string(s.iconData) != string(assets.Icon) {
		t.Error("iconData should match embedded Icon asset")
	}
	if string(s.alertIconData) != string(assets.IconAlert) {
		t.Error("alertIconData should match embedded IconAlert asset")
	}
}

func TestSystray_Setup(t *testing.T) {
	t.Parallel()

	s := NewSystray()
	err := s.Setup()

	if err != nil {
		t.Errorf("Setup should return nil for embedded assets, got: %v", err)
	}
}

func TestSystray_SetMenu(t *testing.T) {
	t.Parallel()

	s := NewSystray()
	menu := &Menu{}

	s.SetMenu(menu)

	if s.menu != menu {
		t.Error("SetMenu did not set menu")
	}
}

func TestSystray_SetAlertState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		hasUnread bool
		wantAlert bool
	}{
		{
			name:      "true shows alert icon",
			hasUnread: true,
			wantAlert: true,
		},
		{
			name:      "false shows normal icon",
			hasUnread: false,
			wantAlert: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewSystray()

			s.SetAlertState(tt.hasUnread)

			if tt.wantAlert {
				if string(s.currentIconData) != string(s.alertIconData) {
					t.Error("SetAlertState(true) should set alert icon as current")
				}
			} else {
				if string(s.currentIconData) != string(s.iconData) {
					t.Error("SetAlertState(false) should set normal icon as current")
				}
			}
		})
	}
}
