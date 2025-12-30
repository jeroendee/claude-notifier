package app

import (
	"os"
	"testing"
)

func TestLoadConfig_DefaultPort(t *testing.T) {
	t.Parallel()

	// Ensure env var is not set
	os.Unsetenv("CLAUDE_NOTIFIER_PORT")

	cfg := LoadConfig()

	if cfg.Port != 19199 {
		t.Errorf("LoadConfig().Port = %d, want 19199", cfg.Port)
	}
}

func TestLoadConfig_CustomPort(t *testing.T) {
	os.Setenv("CLAUDE_NOTIFIER_PORT", "8080")
	defer os.Unsetenv("CLAUDE_NOTIFIER_PORT")

	cfg := LoadConfig()

	if cfg.Port != 8080 {
		t.Errorf("LoadConfig().Port = %d, want 8080", cfg.Port)
	}
}

func TestLoadConfig_InvalidPort(t *testing.T) {
	os.Setenv("CLAUDE_NOTIFIER_PORT", "invalid")
	defer os.Unsetenv("CLAUDE_NOTIFIER_PORT")

	cfg := LoadConfig()

	if cfg.Port != 19199 {
		t.Errorf("LoadConfig().Port = %d, want 19199 (default on invalid)", cfg.Port)
	}
}

func TestLoadConfig_DefaultSoundFile(t *testing.T) {
	t.Parallel()

	os.Unsetenv("CLAUDE_NOTIFIER_SOUND")

	cfg := LoadConfig()

	want := "/System/Library/Sounds/Glass.aiff"
	if cfg.SoundFile != want {
		t.Errorf("LoadConfig().SoundFile = %q, want %q", cfg.SoundFile, want)
	}
}

func TestLoadConfig_CustomSoundFile(t *testing.T) {
	os.Setenv("CLAUDE_NOTIFIER_SOUND", "/custom/sound.wav")
	defer os.Unsetenv("CLAUDE_NOTIFIER_SOUND")

	cfg := LoadConfig()

	want := "/custom/sound.wav"
	if cfg.SoundFile != want {
		t.Errorf("LoadConfig().SoundFile = %q, want %q", cfg.SoundFile, want)
	}
}

func TestLoadConfig_DefaultMaxHistory(t *testing.T) {
	t.Parallel()

	os.Unsetenv("CLAUDE_NOTIFIER_MAX_HISTORY")

	cfg := LoadConfig()

	if cfg.MaxHistory != 50 {
		t.Errorf("LoadConfig().MaxHistory = %d, want 50", cfg.MaxHistory)
	}
}

func TestLoadConfig_CustomMaxHistory(t *testing.T) {
	os.Setenv("CLAUDE_NOTIFIER_MAX_HISTORY", "100")
	defer os.Unsetenv("CLAUDE_NOTIFIER_MAX_HISTORY")

	cfg := LoadConfig()

	if cfg.MaxHistory != 100 {
		t.Errorf("LoadConfig().MaxHistory = %d, want 100", cfg.MaxHistory)
	}
}

func TestLoadConfig_InvalidMaxHistory(t *testing.T) {
	os.Setenv("CLAUDE_NOTIFIER_MAX_HISTORY", "invalid")
	defer os.Unsetenv("CLAUDE_NOTIFIER_MAX_HISTORY")

	cfg := LoadConfig()

	if cfg.MaxHistory != 50 {
		t.Errorf("LoadConfig().MaxHistory = %d, want 50 (default on invalid)", cfg.MaxHistory)
	}
}
