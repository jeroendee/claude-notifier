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

	// Play() should return almost immediately (< 100ms), not wait for command
	if elapsed > 100*time.Millisecond {
		t.Errorf("Play() took %v, expected < 100ms (should be non-blocking)", elapsed)
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
