// ABOUTME: Notifier command - macOS menu bar notification app for Claude Code.
// ABOUTME: This is the main entry point for the notifier application.
package main

import (
	"fmt"
	"os"

	"github.com/jeroendee/claude-notifier/internal/app"
	"github.com/jeroendee/claude-notifier/internal/notification"
	"github.com/jeroendee/claude-notifier/internal/server"
	"github.com/jeroendee/claude-notifier/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := app.LoadConfig()

	store := notification.NewStore()
	soundPlayer := notification.NewSoundPlayer(cfg.SoundFile)
	srv := server.NewServer(cfg.Port, store, soundPlayer)

	tray := ui.NewSystray("assets/icon.png")
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
