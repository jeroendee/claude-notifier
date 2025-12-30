// ABOUTME: Notifier command - macOS menu bar notification app for Claude Code.
// ABOUTME: This is the main entry point for the notifier application.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jeroendee/claude-notifier/internal/app"
	"github.com/jeroendee/claude-notifier/internal/notification"
	"github.com/jeroendee/claude-notifier/internal/server"
	"github.com/jeroendee/claude-notifier/internal/ui"
	"github.com/jeroendee/claude-notifier/internal/version"
)

func main() {
	if handleVersionFlag(os.Stdout) {
		return
	}
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// handleVersionFlag parses --version/-v flags and prints version info if set.
func handleVersionFlag(w io.Writer) bool {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&showVersion, "v", false, "print version and exit (shorthand)")
	flag.Parse()

	if showVersion {
		fmt.Fprintln(w, version.Get())
		return true
	}
	return false
}

func run() error {
	cfg := app.LoadConfig()

	store := notification.NewStore()
	soundPlayer := notification.NewSoundPlayer(cfg.SoundFile)
	srv := server.NewServer(cfg.Port, store, soundPlayer)

	tray := ui.NewSystray()
	menu := ui.NewMenu(tray, store)
	tray.SetMenu(menu)

	application := app.New(cfg).
		SetStore(store).
		SetServer(srv).
		SetSystray(tray).
		SetMenu(menu)

	fmt.Printf("Claude notifier starting on port %d...\n", cfg.Port)
	return application.Run()
}
