package notification

import (
	"os/exec"
	"testing"
	"time"
)

func TestNewSoundPlayer(t *testing.T) {
	tests := []struct {
		name      string
		soundPath string
		wantPath  string
	}{
		{
			name:      "empty path defaults to Glass.aiff",
			soundPath: "",
			wantPath:  "/System/Library/Sounds/Glass.aiff",
		},
		{
			name:      "custom path uses provided path",
			soundPath: "/custom/sound.aiff",
			wantPath:  "/custom/sound.aiff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			player := NewSoundPlayer(tt.soundPath)

			if player.soundPath != tt.wantPath {
				t.Errorf("NewSoundPlayer() soundPath = %q, want %q", player.soundPath, tt.wantPath)
			}
		})
	}
}

func TestSoundPlayer_Play(t *testing.T) {
	tests := []struct {
		name      string
		soundPath string
		wantCmd   string
		wantArgs  []string
	}{
		{
			name:      "executes afplay with default sound path",
			soundPath: "",
			wantCmd:   "afplay",
			wantArgs:  []string{"/System/Library/Sounds/Glass.aiff"},
		},
		{
			name:      "executes afplay with custom sound path",
			soundPath: "/custom/sound.aiff",
			wantCmd:   "afplay",
			wantArgs:  []string{"/custom/sound.aiff"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedCmd string
			var capturedArgs []string

			// Mock execCommand to capture what would be executed
			player := NewSoundPlayer(tt.soundPath)
			player.execCommand = func(name string, args ...string) *exec.Cmd {
				capturedCmd = name
				capturedArgs = args
				// Return a command that does nothing (echo)
				return exec.Command("echo")
			}

			err := player.Play()

			if err != nil {
				t.Errorf("Play() error = %v, want nil", err)
			}

			if capturedCmd != tt.wantCmd {
				t.Errorf("Play() command = %q, want %q", capturedCmd, tt.wantCmd)
			}

			if len(capturedArgs) != len(tt.wantArgs) {
				t.Errorf("Play() args = %v, want %v", capturedArgs, tt.wantArgs)
			} else {
				for i, arg := range capturedArgs {
					if arg != tt.wantArgs[i] {
						t.Errorf("Play() args[%d] = %q, want %q", i, arg, tt.wantArgs[i])
					}
				}
			}
		})
	}
}

func TestSoundPlayer_Play_NonBlocking(t *testing.T) {
	t.Parallel()

	player := NewSoundPlayer("")
	// Use sleep command that takes 2 seconds
	player.execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("sleep", "2")
	}

	start := time.Now()
	err := player.Play()
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Play() error = %v, want nil", err)
	}

	// Play() should return quickly (< 500ms), not wait for the 2s command
	// Using 500ms threshold to account for CI runner variability
	if elapsed > 500*time.Millisecond {
		t.Errorf("Play() took %v, expected < 500ms (should be non-blocking)", elapsed)
	}
}

func TestSoundPlayer_Play_ReturnsErrorOnStartFailure(t *testing.T) {
	t.Parallel()

	player := NewSoundPlayer("")
	// Use a non-existent command to trigger Start() failure
	player.execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("/nonexistent/command/that/does/not/exist")
	}

	err := player.Play()

	if err == nil {
		t.Error("Play() error = nil, want error for failed command start")
	}
}

func TestSoundPlayer_IsMuted_DefaultFalse(t *testing.T) {
	t.Parallel()

	player := NewSoundPlayer("")

	if player.IsMuted() {
		t.Error("IsMuted() = true, want false for new player")
	}
}

func TestSoundPlayer_SetMuted_True(t *testing.T) {
	t.Parallel()

	player := NewSoundPlayer("")

	player.SetMuted(true)

	if !player.IsMuted() {
		t.Error("IsMuted() = false after SetMuted(true), want true")
	}
}

func TestSoundPlayer_SetMuted_False(t *testing.T) {
	t.Parallel()

	player := NewSoundPlayer("")
	player.SetMuted(true)

	player.SetMuted(false)

	if player.IsMuted() {
		t.Error("IsMuted() = true after SetMuted(false), want false")
	}
}

func TestSoundPlayer_Play_SkipsWhenMuted(t *testing.T) {
	t.Parallel()

	commandExecuted := false
	player := NewSoundPlayer("")
	player.execCommand = func(name string, args ...string) *exec.Cmd {
		commandExecuted = true
		return exec.Command("echo")
	}

	player.SetMuted(true)
	err := player.Play()

	if err != nil {
		t.Errorf("Play() error = %v, want nil when muted", err)
	}

	if commandExecuted {
		t.Error("Play() executed command when muted, want no execution")
	}
}

func TestSoundPlayer_Play_ExecutesWhenNotMuted(t *testing.T) {
	t.Parallel()

	commandExecuted := false
	player := NewSoundPlayer("")
	player.execCommand = func(name string, args ...string) *exec.Cmd {
		commandExecuted = true
		return exec.Command("echo")
	}

	player.SetMuted(false)
	err := player.Play()

	if err != nil {
		t.Errorf("Play() error = %v, want nil", err)
	}

	if !commandExecuted {
		t.Error("Play() did not execute command when not muted, want execution")
	}
}

func TestSoundPlayer_MuteState_ThreadSafe(t *testing.T) {
	t.Parallel()

	player := NewSoundPlayer("")
	player.execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("echo")
	}

	done := make(chan bool)
	const goroutines = 100

	// Concurrent writers
	for i := 0; i < goroutines; i++ {
		go func(mute bool) {
			player.SetMuted(mute)
			done <- true
		}(i%2 == 0)
	}

	// Concurrent readers
	for i := 0; i < goroutines; i++ {
		go func() {
			_ = player.IsMuted()
			done <- true
		}()
	}

	// Concurrent Play calls
	for i := 0; i < goroutines; i++ {
		go func() {
			_ = player.Play()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < goroutines*3; i++ {
		<-done
	}
}
