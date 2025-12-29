package notification

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const maxItems = 50

// Notification represents a single notification with read status.
type Notification struct {
	ID        string
	Message   string
	Timestamp time.Time
	Read      bool
}

// Store manages notifications with thread-safe operations and FIFO eviction.
type Store struct {
	mu            sync.RWMutex
	notifications []Notification
	onChange      func()
}

// NewStore creates a new notification store.
func NewStore() *Store {
	return &Store{
		notifications: make([]Notification, 0),
	}
}

// List returns a copy of all notifications.
func (s *Store) List() []Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Notification, len(s.notifications))
	copy(result, s.notifications)
	return result
}

// UnreadCount returns the count of unread notifications.
func (s *Store) UnreadCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, n := range s.notifications {
		if !n.Read {
			count++
		}
	}
	return count
}

// Add creates a new notification with generated ID and current timestamp.
// Evicts oldest notification if store exceeds maxItems.
func (s *Store) Add(message string) {
	s.mu.Lock()
	n := Notification{
		ID:        generateID(),
		Message:   message,
		Timestamp: time.Now(),
		Read:      false,
	}
	s.notifications = append(s.notifications, n)

	// FIFO eviction if exceeding max
	if len(s.notifications) > maxItems {
		s.notifications = s.notifications[1:]
	}

	onChange := s.onChange
	s.mu.Unlock()

	if onChange != nil {
		onChange()
	}
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// SetOnChange sets the callback function called when notifications change.
func (s *Store) SetOnChange(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onChange = fn
}

// MarkRead marks a specific notification as read by ID.
func (s *Store) MarkRead(id string) {
	s.mu.Lock()
	found := false
	for i := range s.notifications {
		if s.notifications[i].ID == id {
			s.notifications[i].Read = true
			found = true
			break
		}
	}
	onChange := s.onChange
	s.mu.Unlock()

	if found && onChange != nil {
		onChange()
	}
}

// MarkAllRead marks all notifications as read.
func (s *Store) MarkAllRead() {
	s.mu.Lock()
	changed := false
	for i := range s.notifications {
		if !s.notifications[i].Read {
			s.notifications[i].Read = true
			changed = true
		}
	}
	onChange := s.onChange
	s.mu.Unlock()

	if changed && onChange != nil {
		onChange()
	}
}

// Clear removes all notifications from the store.
func (s *Store) Clear() {
	s.mu.Lock()
	hadNotifications := len(s.notifications) > 0
	s.notifications = make([]Notification, 0)
	onChange := s.onChange
	s.mu.Unlock()

	if hadNotifications && onChange != nil {
		onChange()
	}
}
