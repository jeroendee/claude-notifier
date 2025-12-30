package app

import (
	"os"
	"os/signal"
	"syscall"
)

// Store defines the interface for notification storage.
type Store interface {
	SetOnChange(fn func())
	UnreadCount() int
}

// Server defines the interface for the HTTP server.
type Server interface {
	Start() error
	Stop() error
}

// Systray defines the interface for system tray management.
type Systray interface {
	Setup() error
	Run()
	SetAlertState(hasUnread bool)
}

// Menu defines the interface for the system tray menu.
type Menu interface {
	Refresh()
}

// App coordinates all application components.
type App struct {
	config  *Config
	store   Store
	server  Server
	systray Systray
	menu    Menu
}

// New creates a new App with the given configuration.
func New(config *Config) *App {
	return &App{
		config: config,
	}
}

// Run starts all components and blocks until shutdown signal.
func (a *App) Run() error {
	// Setup systray first
	if a.systray != nil {
		if err := a.systray.Setup(); err != nil {
			return err
		}
	}

	// Set up onChange callback
	a.setupOnChange()

	// Set initial icon state based on existing notifications
	if a.systray != nil && a.store != nil {
		a.systray.SetAlertState(a.store.UnreadCount() > 0)
	}

	// Start HTTP server
	if a.server != nil {
		if err := a.server.Start(); err != nil {
			return err
		}
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Run systray or wait for signal
	if a.systray != nil {
		// Run systray in background and handle signals
		go func() {
			<-sigChan
			a.Shutdown()
		}()
		a.systray.Run()
		a.Shutdown()
	} else {
		<-sigChan
		a.Shutdown()
	}

	return nil
}

// setupOnChange configures the store's onChange callback to refresh the menu
// and update the systray icon alert state.
func (a *App) setupOnChange() {
	if a.store != nil && a.menu != nil {
		a.store.SetOnChange(func() {
			a.menu.Refresh()
			if a.systray != nil {
				a.systray.SetAlertState(a.store.UnreadCount() > 0)
			}
		})
	}
}

// Shutdown gracefully stops all components.
func (a *App) Shutdown() {
	if a.server != nil {
		a.server.Stop()
	}
}

// SetStore sets the notification store and returns the App for chaining.
func (a *App) SetStore(store Store) *App {
	a.store = store
	return a
}

// SetServer sets the HTTP server and returns the App for chaining.
func (a *App) SetServer(server Server) *App {
	a.server = server
	return a
}

// SetSystray sets the system tray and returns the App for chaining.
func (a *App) SetSystray(systray Systray) *App {
	a.systray = systray
	return a
}

// SetMenu sets the menu and returns the App for chaining.
func (a *App) SetMenu(menu Menu) *App {
	a.menu = menu
	return a
}
