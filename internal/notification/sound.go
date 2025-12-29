package notification

import "os/exec"

const defaultSoundPath = "/System/Library/Sounds/Glass.aiff"

// SoundPlayer plays notification sounds using macOS afplay command.
type SoundPlayer struct {
	soundPath   string
	execCommand func(name string, arg ...string) *exec.Cmd
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
func (p *SoundPlayer) Play() error {
	cmd := p.execCommand("afplay", p.soundPath)
	return cmd.Start()
}
