package app

import (
	"os"
	"strings"
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

func TestLoadConfig_PortValidation(t *testing.T) {
	tests := []struct {
		name      string
		portValue string
		wantErr   bool
		errSubstr string
	}{
		{
			name:      "port zero rejected",
			portValue: "0",
			wantErr:   true,
			errSubstr: "port must be between 1 and 65535",
		},
		{
			name:      "negative port rejected",
			portValue: "-1",
			wantErr:   true,
			errSubstr: "port must be between 1 and 65535",
		},
		{
			name:      "port above max rejected",
			portValue: "65536",
			wantErr:   true,
			errSubstr: "port must be between 1 and 65535",
		},
		{
			name:      "port at minimum accepted",
			portValue: "1",
			wantErr:   false,
		},
		{
			name:      "port at maximum accepted",
			portValue: "65535",
			wantErr:   false,
		},
		{
			name:      "port in valid range accepted",
			portValue: "8080",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("CLAUDE_NOTIFIER_PORT", tt.portValue)
			defer os.Unsetenv("CLAUDE_NOTIFIER_PORT")

			_, err := LoadConfigWithValidation()

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadConfigWithValidation() error = nil, want error containing %q", tt.errSubstr)
					return
				}
				if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("LoadConfigWithValidation() error = %q, want error containing %q", err.Error(), tt.errSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("LoadConfigWithValidation() unexpected error = %v", err)
				}
			}
		})
	}
}
