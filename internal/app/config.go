package app

import (
	"errors"
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

// LoadConfigWithValidation loads configuration and validates all values.
// Returns an error if any configuration value is invalid.
func LoadConfigWithValidation() (*Config, error) {
	cfg := &Config{
		Port:      19199,
		SoundFile: "/System/Library/Sounds/Glass.aiff",
	}

	if portStr := os.Getenv("CLAUDE_NOTIFIER_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, errors.New("invalid port: must be a number")
		}
		cfg.Port = port
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, errors.New("invalid port: port must be between 1 and 65535")
	}

	if soundFile := os.Getenv("CLAUDE_NOTIFIER_SOUND"); soundFile != "" {
		cfg.SoundFile = soundFile
	}

	return cfg, nil
}
