package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jeroendee/claude-notifier/internal/app"
	"github.com/jeroendee/claude-notifier/internal/notification"
	"github.com/jeroendee/claude-notifier/internal/server"
	"github.com/jeroendee/claude-notifier/internal/transcript"
	"github.com/jeroendee/claude-notifier/internal/ui"
	"github.com/jeroendee/claude-notifier/internal/version"
)

func main() {
	if handled, err := handleSummarySubcommand(os.Args[1:], os.Stdout); handled {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if handleVersionFlag(os.Stdout) {
		return
	}
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// handleSummarySubcommand handles the 'summary' subcommand.
// Returns (true, nil) if handled successfully, (true, error) if handled with error,
// or (false, nil) if args don't contain the summary subcommand.
func handleSummarySubcommand(args []string, w io.Writer) (bool, error) {
	if len(args) == 0 || args[0] != "summary" {
		return false, nil
	}

	if len(args) < 2 {
		return true, errors.New("summary: missing transcript path argument")
	}

	path := args[1]
	messages, err := transcript.ParseTranscript(path)
	if err != nil {
		return true, fmt.Errorf("summary: %w", err)
	}

	heading := transcript.ExtractLastHeading(messages)
	if heading != "" {
		fmt.Fprintln(w, heading)
	}

	return true, nil
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

	fmt.Printf("claude-notifier starting on port %d...\n", cfg.Port)
	return application.Run()
}
