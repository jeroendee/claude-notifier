package app

import (
	"os"
	"strconv"
)

// Config holds application configuration.
type Config struct {
	Port      int
	SoundFile string
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() *Config {
	cfg := &Config{
		Port:      19199,
		SoundFile: "/System/Library/Sounds/Glass.aiff",
	}

	if portStr := os.Getenv("CLAUDE_NOTIFIER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	if soundFile := os.Getenv("CLAUDE_NOTIFIER_SOUND"); soundFile != "" {
		cfg.SoundFile = soundFile
	}

	return cfg
}
