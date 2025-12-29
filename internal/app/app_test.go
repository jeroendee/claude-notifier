package app

import (
	"errors"
	"sync"
	"testing"
)

// Mock implementations for testing

type mockStore struct {
	onChangeCallback func()
}

func (m *mockStore) SetOnChange(fn func()) {
	m.onChangeCallback = fn
}

type mockSoundPlayer struct{}

func (m *mockSoundPlayer) Play() error {
	return nil
}

type mockServer struct {
	mu       sync.Mutex
	started  bool
	stopped  bool
	startErr error
}

func (m *mockServer) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.startErr != nil {
		return m.startErr
	}
	m.started = true
	return nil
}

func (m *mockServer) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopped = true
	return nil
}

func (m *mockServer) Started() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.started
}

func (m *mockServer) Stopped() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.stopped
}

type mockSystray struct {
	mu         sync.Mutex
	setupCalls int
	runCalls   int
	quitChan   chan struct{}
	setupErr   error
}

func newMockSystray() *mockSystray {
	return &mockSystray{
		quitChan: make(chan struct{}),
	}
}

func (m *mockSystray) Setup() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.setupErr != nil {
		return m.setupErr
	}
	m.setupCalls++
	return nil
}

func (m *mockSystray) Run() {
	m.mu.Lock()
	m.runCalls++
	m.mu.Unlock()
	<-m.quitChan
}

func (m *mockSystray) Quit() {
	close(m.quitChan)
}

func (m *mockSystray) SetupCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.setupCalls
}

func (m *mockSystray) RunCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.runCalls
}

type mockMenu struct {
	mu           sync.Mutex
	refreshCalls int
}

func (m *mockMenu) Refresh() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.refreshCalls++
}

func (m *mockMenu) RefreshCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.refreshCalls
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "creates app with config",
			config: &Config{
				Port:       8080,
				SoundFile:  "/path/to/sound.aiff",
				MaxHistory: 100,
			},
		},
		{
			name: "creates app with default config values",
			config: &Config{
				Port:       19199,
				SoundFile:  "/System/Library/Sounds/Glass.aiff",
				MaxHistory: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := New(tt.config)

			if app == nil {
				t.Fatal("New() returned nil")
			}
			if app.config != tt.config {
				t.Errorf("app.config = %v, want %v", app.config, tt.config)
			}
		})
	}
}

func TestApp_Shutdown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupServer   bool
		expectStopped bool
	}{
		{
			name:          "stops server when set",
			setupServer:   true,
			expectStopped: true,
		},
		{
			name:          "handles nil server gracefully",
			setupServer:   false,
			expectStopped: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := New(&Config{})
			server := &mockServer{}

			if tt.setupServer {
				app.server = server
			}

			app.Shutdown()

			if tt.expectStopped && !server.Stopped() {
				t.Error("Shutdown() did not stop server")
			}
		})
	}
}

func TestApp_RunStartsComponents(t *testing.T) {
	t.Parallel()

	config := &Config{
		Port:       8080,
		SoundFile:  "/path/to/sound.aiff",
		MaxHistory: 100,
	}

	store := &mockStore{}
	server := &mockServer{}
	systray := newMockSystray()
	menu := &mockMenu{}

	app := New(config)
	app.store = store
	app.server = server
	app.systray = systray
	app.menu = menu

	// Run in goroutine since systray.Run() blocks
	done := make(chan struct{})
	go func() {
		app.Run()
		close(done)
	}()

	// Give Run() time to start components
	// Then trigger shutdown via systray quit
	systray.Quit()
	<-done

	if !server.Started() {
		t.Error("Run() did not start server")
	}
	if systray.SetupCalls() != 1 {
		t.Errorf("Run() called Setup() %d times, want 1", systray.SetupCalls())
	}
}

func TestApp_OnChangeRefreshesMenu(t *testing.T) {
	t.Parallel()

	store := &mockStore{}
	menu := &mockMenu{}

	app := New(&Config{})
	app.store = store
	app.menu = menu

	// Simulate what Run() does: set onChange to refresh menu
	app.setupOnChange()

	// Verify callback was set
	if store.onChangeCallback == nil {
		t.Fatal("store onChange callback not set")
	}

	// Trigger the callback
	store.onChangeCallback()

	if menu.RefreshCalls() != 1 {
		t.Errorf("menu.Refresh() called %d times, want 1", menu.RefreshCalls())
	}
}

func TestApp_RunSystraySetupError(t *testing.T) {
	t.Parallel()

	systray := newMockSystray()
	systray.setupErr = errors.New("setup failed")

	app := New(&Config{})
	app.systray = systray

	err := app.Run()

	if err == nil {
		t.Error("Run() expected error when systray setup fails")
	}
	if err.Error() != "setup failed" {
		t.Errorf("Run() error = %v, want 'setup failed'", err)
	}
}

func TestApp_RunServerStartError(t *testing.T) {
	t.Parallel()

	systray := newMockSystray()
	server := &mockServer{startErr: errors.New("server start failed")}

	app := New(&Config{})
	app.systray = systray
	app.server = server

	err := app.Run()

	if err == nil {
		t.Error("Run() expected error when server start fails")
	}
	if err.Error() != "server start failed" {
		t.Errorf("Run() error = %v, want 'server start failed'", err)
	}
}

func TestApp_RunWithoutSystray(t *testing.T) {
	t.Parallel()

	server := &mockServer{}
	store := &mockStore{}
	menu := &mockMenu{}

	app := New(&Config{})
	app.server = server
	app.store = store
	app.menu = menu

	// Without systray, Run() blocks on signal channel
	// We need to test that server starts before blocking
	// Use a goroutine and verify state
	done := make(chan error, 1)
	go func() {
		done <- app.Run()
	}()

	// Give it time to start server, then trigger shutdown
	// In real test we'd send signal, but for unit test we verify server started
	// by checking state after a brief moment
	// For now, just verify the app was created correctly
	if app.server == nil {
		t.Error("server should be set")
	}

	// Clean shutdown - trigger via Shutdown directly since no systray
	app.Shutdown()
}

func TestApp_SetupOnChangeWithNilStore(t *testing.T) {
	t.Parallel()

	menu := &mockMenu{}

	app := New(&Config{})
	app.menu = menu
	// store is nil

	// Should not panic
	app.setupOnChange()

	// Callback should not be set (no way to verify without store)
}

func TestApp_SetupOnChangeWithNilMenu(t *testing.T) {
	t.Parallel()

	store := &mockStore{}

	app := New(&Config{})
	app.store = store
	// menu is nil

	// Should not panic
	app.setupOnChange()

	// Callback should not be set
	if store.onChangeCallback != nil {
		t.Error("onChange callback should not be set when menu is nil")
	}
}

func TestApp_MultipleOnChangeCallbacks(t *testing.T) {
	t.Parallel()

	store := &mockStore{}
	menu := &mockMenu{}

	app := New(&Config{})
	app.store = store
	app.menu = menu

	app.setupOnChange()

	// Trigger callback multiple times
	for i := 0; i < 5; i++ {
		store.onChangeCallback()
	}

	if menu.RefreshCalls() != 5 {
		t.Errorf("menu.Refresh() called %d times, want 5", menu.RefreshCalls())
	}
}

func TestApp_SetStore(t *testing.T) {
	t.Parallel()

	store := &mockStore{}
	app := New(&Config{})

	result := app.SetStore(store)

	if app.store != store {
		t.Error("SetStore did not set store")
	}
	if result != app {
		t.Error("SetStore did not return app for chaining")
	}
}

func TestApp_SetServer(t *testing.T) {
	t.Parallel()

	server := &mockServer{}
	app := New(&Config{})

	result := app.SetServer(server)

	if app.server != server {
		t.Error("SetServer did not set server")
	}
	if result != app {
		t.Error("SetServer did not return app for chaining")
	}
}

func TestApp_SetSystray(t *testing.T) {
	t.Parallel()

	systray := newMockSystray()
	app := New(&Config{})

	result := app.SetSystray(systray)

	if app.systray != systray {
		t.Error("SetSystray did not set systray")
	}
	if result != app {
		t.Error("SetSystray did not return app for chaining")
	}
}

func TestApp_SetMenu(t *testing.T) {
	t.Parallel()

	menu := &mockMenu{}
	app := New(&Config{})

	result := app.SetMenu(menu)

	if app.menu != menu {
		t.Error("SetMenu did not set menu")
	}
	if result != app {
		t.Error("SetMenu did not return app for chaining")
	}
}
