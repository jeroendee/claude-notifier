package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jeroendee/claude-notifier/internal/notification"
)

// mockSoundPlayer tracks Play calls for testing.
type mockSoundPlayer struct {
	playCalled bool
	playErr    error
}

func (m *mockSoundPlayer) Play() error {
	m.playCalled = true
	return m.playErr
}

// mockNotificationStore implements NotificationStore for testing.
type mockNotificationStore struct {
	addCalled   bool
	addMessage  string
	clearCalled bool
}

func (m *mockNotificationStore) Add(message string) {
	m.addCalled = true
	m.addMessage = message
}

func (m *mockNotificationStore) Clear() {
	m.clearCalled = true
}

func TestNewServer(t *testing.T) {
	t.Parallel()

	store := notification.NewStore()
	player := &mockSoundPlayer{}

	srv := NewServer(8080, store, player, nil)

	if srv == nil {
		t.Fatal("NewServer() returned nil")
	}
	if srv.port != 8080 {
		t.Errorf("port = %d, want 8080", srv.port)
	}
	if srv.store != store {
		t.Error("store not set correctly")
	}
	if srv.player != player {
		t.Error("player not set correctly")
	}
}

func TestNewServerWithMockStore(t *testing.T) {
	t.Parallel()

	store := &mockNotificationStore{}
	player := &mockSoundPlayer{}

	srv := NewServer(8080, store, player, nil)

	if srv == nil {
		t.Fatal("NewServer() returned nil")
	}
	if srv.port != 8080 {
		t.Errorf("port = %d, want 8080", srv.port)
	}
}

func TestNotifyHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantPlayCalled bool
		wantStoreCount int
	}{
		{
			name:           "valid notification",
			method:         http.MethodPost,
			body:           `{"message": "test notification"}`,
			wantStatus:     http.StatusCreated,
			wantPlayCalled: true,
			wantStoreCount: 1,
		},
		{
			name:           "invalid JSON",
			method:         http.MethodPost,
			body:           `{invalid`,
			wantStatus:     http.StatusBadRequest,
			wantPlayCalled: false,
			wantStoreCount: 0,
		},
		{
			name:           "empty message",
			method:         http.MethodPost,
			body:           `{"message": ""}`,
			wantStatus:     http.StatusBadRequest,
			wantPlayCalled: false,
			wantStoreCount: 0,
		},
		{
			name:           "missing message field",
			method:         http.MethodPost,
			body:           `{}`,
			wantStatus:     http.StatusBadRequest,
			wantPlayCalled: false,
			wantStoreCount: 0,
		},
		{
			name:           "wrong method GET",
			method:         http.MethodGet,
			body:           "",
			wantStatus:     http.StatusMethodNotAllowed,
			wantPlayCalled: false,
			wantStoreCount: 0,
		},
		{
			name:           "wrong method PUT",
			method:         http.MethodPut,
			body:           `{"message": "test"}`,
			wantStatus:     http.StatusMethodNotAllowed,
			wantPlayCalled: false,
			wantStoreCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := notification.NewStore()
			player := &mockSoundPlayer{}
			srv := NewServer(8080, store, player, nil)

			req := httptest.NewRequest(tt.method, "/notify", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			srv.handleNotify(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if player.playCalled != tt.wantPlayCalled {
				t.Errorf("playCalled = %v, want %v", player.playCalled, tt.wantPlayCalled)
			}
			if len(store.List()) != tt.wantStoreCount {
				t.Errorf("store count = %d, want %d", len(store.List()), tt.wantStoreCount)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "GET returns ok",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok"}`,
		},
		{
			name:       "POST not allowed",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "",
		},
		{
			name:       "PUT not allowed",
			method:     http.MethodPut,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := notification.NewStore()
			player := &mockSoundPlayer{}
			srv := NewServer(8080, store, player, nil)

			req := httptest.NewRequest(tt.method, "/health", nil)
			rec := httptest.NewRecorder()

			srv.handleHealth(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantBody != "" {
				var got, want map[string]string
				json.Unmarshal(rec.Body.Bytes(), &got)
				json.Unmarshal([]byte(tt.wantBody), &want)
				if got["status"] != want["status"] {
					t.Errorf("body status = %q, want %q", got["status"], want["status"])
				}
			}
		})
	}
}

func TestClearHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{
			name:       "POST clears notifications",
			method:     http.MethodPost,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "GET not allowed",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "PUT not allowed",
			method:     http.MethodPut,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := notification.NewStore()
			player := &mockSoundPlayer{}
			srv := NewServer(8080, store, player, nil)

			// Add a notification first
			store.Add("test notification")
			if len(store.List()) != 1 {
				t.Fatal("setup failed: notification not added")
			}

			req := httptest.NewRequest(tt.method, "/clear", nil)
			rec := httptest.NewRecorder()

			srv.handleClear(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			// Verify clear happened only for POST
			if tt.method == http.MethodPost && len(store.List()) != 0 {
				t.Errorf("store not cleared, count = %d", len(store.List()))
			}
			if tt.method != http.MethodPost && len(store.List()) != 1 {
				t.Errorf("store should not be cleared for method %s", tt.method)
			}
		})
	}
}

func TestServerStartStop(t *testing.T) {
	t.Parallel()

	store := notification.NewStore()
	player := &mockSoundPlayer{}
	srv := NewServer(0, store, player, nil) // port 0 = auto-assign

	if err := srv.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Server should be running
	if srv.server == nil {
		t.Error("server.server is nil after Start()")
	}

	if err := srv.Stop(); err != nil {
		t.Errorf("Stop() error = %v", err)
	}
}

func TestServerRouting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{
			name:       "notify endpoint",
			method:     http.MethodPost,
			path:       "/notify",
			body:       `{"message": "test"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "health endpoint",
			method:     http.MethodGet,
			path:       "/health",
			body:       "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "clear endpoint",
			method:     http.MethodPost,
			path:       "/clear",
			body:       "",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "unknown endpoint",
			method:     http.MethodGet,
			path:       "/unknown",
			body:       "",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Not parallel - uses shared store
			store := notification.NewStore()
			player := &mockSoundPlayer{}
			srv := NewServer(8080, store, player, nil)

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()

			srv.mux().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestNotifyHandler_LogsSoundPlayerErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		playErr     error
		wantLogged  bool
		wantContain string
	}{
		{
			name:        "logs error when Play fails",
			playErr:     errors.New("sound file not found"),
			wantLogged:  true,
			wantContain: "sound file not found",
		},
		{
			name:       "no error logged when Play succeeds",
			playErr:    nil,
			wantLogged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := notification.NewStore()
			player := &mockSoundPlayer{playErr: tt.playErr}

			// Capture log output
			var logBuf strings.Builder
			logger := log.New(&logBuf, "", 0)

			srv := NewServer(8080, store, player, logger)

			req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{"message": "test"}`))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			srv.handleNotify(rec, req)

			// Notification should still succeed even if sound fails
			if rec.Code != http.StatusCreated {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
			}

			// Verify notification was stored
			if len(store.List()) != 1 {
				t.Errorf("store count = %d, want 1", len(store.List()))
			}

			logOutput := logBuf.String()
			if tt.wantLogged {
				if !strings.Contains(logOutput, tt.wantContain) {
					t.Errorf("log output = %q, want to contain %q", logOutput, tt.wantContain)
				}
			} else {
				if logOutput != "" {
					t.Errorf("log output = %q, want empty", logOutput)
				}
			}
		})
	}
}
