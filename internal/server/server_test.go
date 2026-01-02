package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	srv := NewServer(8080, store, player)

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

	srv := NewServer(8080, store, player)

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
			srv := NewServer(8080, store, player)

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
			srv := NewServer(8080, store, player)

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
			srv := NewServer(8080, store, player)

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
	srv := NewServer(0, store, player) // port 0 = auto-assign

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
			srv := NewServer(8080, store, player)

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
