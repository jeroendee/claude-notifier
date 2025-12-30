package app

import (
	"os"
	"strconv"
)

// Config holds application configuration.
type Config struct {
	Port       int
	SoundFile  string
	MaxHistory int
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() *Config {
	cfg := &Config{
		Port:       19199,
		SoundFile:  "/System/Library/Sounds/Glass.aiff",
		MaxHistory: 50,
	}

	if portStr := os.Getenv("CLAUDE_NOTIFIER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	if soundFile := os.Getenv("CLAUDE_NOTIFIER_SOUND"); soundFile != "" {
		cfg.SoundFile = soundFile
	}

	if maxHistStr := os.Getenv("CLAUDE_NOTIFIER_MAX_HISTORY"); maxHistStr != "" {
		if maxHist, err := strconv.Atoi(maxHistStr); err == nil {
			cfg.MaxHistory = maxHist
		}
	}

	return cfg
}
