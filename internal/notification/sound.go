package notification

import (
	"os/exec"
	"sync"
)

const defaultSoundPath = "/System/Library/Sounds/Glass.aiff"

// SoundPlayer plays notification sounds using macOS afplay command.
type SoundPlayer struct {
	soundPath   string
	execCommand func(name string, arg ...string) *exec.Cmd
	muted       bool
	mu          sync.RWMutex
}

// NewSoundPlayer creates a new SoundPlayer. If soundPath is empty,
// defaults to /System/Library/Sounds/Glass.aiff.
func NewSoundPlayer(soundPath string) *SoundPlayer {
	if soundPath == "" {
		soundPath = defaultSoundPath
	}
	return &SoundPlayer{
		soundPath:   soundPath,
		execCommand: exec.Command,
	}
}

// Play runs afplay in background (non-blocking) to play the sound.
// Returns nil immediately if muted.
func (p *SoundPlayer) Play() error {
	if p.IsMuted() {
		return nil
	}
	cmd := p.execCommand("afplay", p.soundPath)
	return cmd.Start()
}

// SetMuted enables or disables sound playback.
func (p *SoundPlayer) SetMuted(muted bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.muted = muted
}

// IsMuted returns true if sound playback is disabled.
func (p *SoundPlayer) IsMuted() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.muted
}
