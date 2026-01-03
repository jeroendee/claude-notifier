package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// NotificationStore defines the interface for storing notifications.
type NotificationStore interface {
	Add(message string)
	Clear()
}

// SoundPlayer defines the interface for playing notification sounds.
type SoundPlayer interface {
	Play() error
}

// NotifyRequest represents the JSON body for POST /notify.
type NotifyRequest struct {
	Message string `json:"message"`
}

// Server handles HTTP requests for notifications.
type Server struct {
	port     int
	store    NotificationStore
	player   SoundPlayer
	logger   *log.Logger
	server   *http.Server
	listener net.Listener
}

// NewServer creates a new notification HTTP server.
// If logger is nil, a default logger writing to stderr is used.
func NewServer(port int, store NotificationStore, player SoundPlayer, logger *log.Logger) *Server {
	if logger == nil {
		logger = log.Default()
	}
	return &Server{
		port:   port,
		store:  store,
		player: player,
		logger: logger,
	}
}

// Start begins listening for HTTP requests. Non-blocking.
func (s *Server) Start() error {
	addr := fmt.Sprintf("localhost:%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}
	s.listener = listener

	s.server = &http.Server{
		Handler:      s.mux(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go s.server.Serve(listener)
	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// mux creates the HTTP router.
func (s *Server) mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/notify", s.handleNotify)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/clear", s.handleClear)
	return mux
}

// handleNotify processes POST /notify requests.
func (s *Server) handleNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	s.store.Add(req.Message)
	if err := s.player.Play(); err != nil {
		s.logger.Printf("sound player error: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
}

// handleHealth processes GET /health requests.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleClear processes POST /clear requests.
func (s *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.store.Clear()
	w.WriteHeader(http.StatusNoContent)
}
