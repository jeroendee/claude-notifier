package notification

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewStore_CreatesEmptyStore(t *testing.T) {
	t.Parallel()

	store := NewStore()

	if store == nil {
		t.Fatal("NewStore() returned nil")
	}

	notifications := store.List()
	if len(notifications) != 0 {
		t.Errorf("NewStore() notifications = %d, want 0", len(notifications))
	}

	if count := store.UnreadCount(); count != 0 {
		t.Errorf("NewStore() UnreadCount() = %d, want 0", count)
	}
}

func TestAdd_CreatesNotificationWithIDAndTimestamp(t *testing.T) {
	t.Parallel()

	store := NewStore()
	before := time.Now()

	store.Add("test message")

	after := time.Now()
	notifications := store.List()

	if len(notifications) != 1 {
		t.Fatalf("Add() notifications count = %d, want 1", len(notifications))
	}

	n := notifications[0]

	if n.ID == "" {
		t.Error("Add() notification ID is empty")
	}

	if n.Message != "test message" {
		t.Errorf("Add() message = %q, want %q", n.Message, "test message")
	}

	if n.Timestamp.Before(before) || n.Timestamp.After(after) {
		t.Errorf("Add() timestamp = %v, want between %v and %v", n.Timestamp, before, after)
	}

	if n.Read {
		t.Error("Add() notification should be unread")
	}
}

func TestAdd_TriggersOnChangeCallback(t *testing.T) {
	t.Parallel()

	store := NewStore()
	callCount := 0
	store.SetOnChange(func() {
		callCount++
	})

	store.Add("message 1")
	store.Add("message 2")

	if callCount != 2 {
		t.Errorf("Add() onChange call count = %d, want 2", callCount)
	}
}

func TestAdd_EvictsOldestWhenExceedingMaxItems(t *testing.T) {
	t.Parallel()

	store := NewStore()

	// Add 51 notifications (exceeds max of 50)
	for i := 0; i < 51; i++ {
		store.Add(fmt.Sprintf("message %d", i))
	}

	notifications := store.List()

	if len(notifications) != 50 {
		t.Fatalf("Add() after overflow, count = %d, want 50", len(notifications))
	}

	// First notification should be "message 1" (message 0 was evicted)
	if notifications[0].Message != "message 1" {
		t.Errorf("Add() oldest notification = %q, want %q", notifications[0].Message, "message 1")
	}

	// Last notification should be "message 50"
	if notifications[49].Message != "message 50" {
		t.Errorf("Add() newest notification = %q, want %q", notifications[49].Message, "message 50")
	}
}

func TestMarkRead_MarksSpecificNotificationAsRead(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Add("message 1")
	store.Add("message 2")

	notifications := store.List()
	targetID := notifications[0].ID

	store.MarkRead(targetID)

	updated := store.List()

	if !updated[0].Read {
		t.Error("MarkRead() notification should be read")
	}

	if updated[1].Read {
		t.Error("MarkRead() other notification should still be unread")
	}

	if store.UnreadCount() != 1 {
		t.Errorf("MarkRead() UnreadCount() = %d, want 1", store.UnreadCount())
	}
}

func TestMarkRead_TriggersOnChangeCallback(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Add("message")

	callCount := 0
	store.SetOnChange(func() {
		callCount++
	})

	notifications := store.List()
	store.MarkRead(notifications[0].ID)

	if callCount != 1 {
		t.Errorf("MarkRead() onChange call count = %d, want 1", callCount)
	}
}

func TestMarkRead_NonExistentIDIsNoOp(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Add("message")

	callCount := 0
	store.SetOnChange(func() {
		callCount++
	})

	store.MarkRead("nonexistent-id")

	if callCount != 0 {
		t.Errorf("MarkRead() with nonexistent ID onChange call count = %d, want 0", callCount)
	}

	if store.UnreadCount() != 1 {
		t.Errorf("MarkRead() with nonexistent ID UnreadCount() = %d, want 1", store.UnreadCount())
	}
}

func TestMarkAllRead_MarksAllNotificationsAsRead(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Add("message 1")
	store.Add("message 2")
	store.Add("message 3")

	store.MarkAllRead()

	notifications := store.List()
	for i, n := range notifications {
		if !n.Read {
			t.Errorf("MarkAllRead() notification[%d] should be read", i)
		}
	}

	if store.UnreadCount() != 0 {
		t.Errorf("MarkAllRead() UnreadCount() = %d, want 0", store.UnreadCount())
	}
}

func TestMarkAllRead_TriggersOnChangeCallback(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Add("message 1")
	store.Add("message 2")

	callCount := 0
	store.SetOnChange(func() {
		callCount++
	})

	store.MarkAllRead()

	if callCount != 1 {
		t.Errorf("MarkAllRead() onChange call count = %d, want 1", callCount)
	}
}

func TestMarkAllRead_EmptyStoreIsNoOp(t *testing.T) {
	t.Parallel()

	store := NewStore()

	callCount := 0
	store.SetOnChange(func() {
		callCount++
	})

	store.MarkAllRead()

	if callCount != 0 {
		t.Errorf("MarkAllRead() on empty store onChange call count = %d, want 0", callCount)
	}
}

func TestClear_RemovesAllNotifications(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Add("message 1")
	store.Add("message 2")
	store.Add("message 3")

	store.Clear()

	notifications := store.List()
	if len(notifications) != 0 {
		t.Errorf("Clear() notifications count = %d, want 0", len(notifications))
	}

	if store.UnreadCount() != 0 {
		t.Errorf("Clear() UnreadCount() = %d, want 0", store.UnreadCount())
	}
}

func TestClear_TriggersOnChangeCallback(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Add("message 1")
	store.Add("message 2")

	callCount := 0
	store.SetOnChange(func() {
		callCount++
	})

	store.Clear()

	if callCount != 1 {
		t.Errorf("Clear() onChange call count = %d, want 1", callCount)
	}
}

func TestClear_EmptyStoreIsNoOp(t *testing.T) {
	t.Parallel()

	store := NewStore()

	callCount := 0
	store.SetOnChange(func() {
		callCount++
	})

	store.Clear()

	if callCount != 0 {
		t.Errorf("Clear() on empty store onChange call count = %d, want 0", callCount)
	}
}

func TestList_ReturnsCopyOfNotifications(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Add("message 1")
	store.Add("message 2")

	// Get list and modify the copy
	list1 := store.List()
	list1[0].Message = "modified"
	list1[0].Read = true

	// Get list again - should be unmodified
	list2 := store.List()

	if list2[0].Message != "message 1" {
		t.Errorf("List() returned reference not copy, message = %q, want %q", list2[0].Message, "message 1")
	}

	if list2[0].Read {
		t.Error("List() returned reference not copy, Read status was modified")
	}
}

func TestUnreadCount_ReturnsCorrectCount(t *testing.T) {
	t.Parallel()

	store := NewStore()

	// Empty store
	if store.UnreadCount() != 0 {
		t.Errorf("UnreadCount() on empty store = %d, want 0", store.UnreadCount())
	}

	// Add notifications
	store.Add("message 1")
	store.Add("message 2")
	store.Add("message 3")

	if store.UnreadCount() != 3 {
		t.Errorf("UnreadCount() after adding 3 = %d, want 3", store.UnreadCount())
	}

	// Mark one as read
	notifications := store.List()
	store.MarkRead(notifications[0].ID)

	if store.UnreadCount() != 2 {
		t.Errorf("UnreadCount() after marking 1 read = %d, want 2", store.UnreadCount())
	}

	// Mark all as read
	store.MarkAllRead()

	if store.UnreadCount() != 0 {
		t.Errorf("UnreadCount() after marking all read = %d, want 0", store.UnreadCount())
	}
}

func TestSetOnChange_ReplacesCallback(t *testing.T) {
	t.Parallel()

	store := NewStore()

	callCount1 := 0
	callCount2 := 0

	store.SetOnChange(func() {
		callCount1++
	})

	store.Add("message 1")

	if callCount1 != 1 {
		t.Errorf("First callback call count = %d, want 1", callCount1)
	}

	// Replace callback
	store.SetOnChange(func() {
		callCount2++
	})

	store.Add("message 2")

	if callCount1 != 1 {
		t.Errorf("First callback call count after replace = %d, want 1", callCount1)
	}

	if callCount2 != 1 {
		t.Errorf("Second callback call count = %d, want 1", callCount2)
	}
}

func TestSetOnChange_NilCallback(t *testing.T) {
	t.Parallel()

	store := NewStore()

	callCount := 0
	store.SetOnChange(func() {
		callCount++
	})

	store.Add("message 1")

	// Set nil callback
	store.SetOnChange(nil)

	store.Add("message 2")

	if callCount != 1 {
		t.Errorf("Callback should not be called after set to nil, call count = %d, want 1", callCount)
	}
}

func TestStore_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	store := NewStore()
	const goroutines = 10
	const opsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 4) // 4 types of operations

	// Concurrent Add
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				store.Add(fmt.Sprintf("message-%d-%d", id, j))
			}
		}(i)
	}

	// Concurrent List
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				_ = store.List()
			}
		}()
	}

	// Concurrent UnreadCount
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				_ = store.UnreadCount()
			}
		}()
	}

	// Concurrent MarkAllRead
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				store.MarkAllRead()
			}
		}()
	}

	wg.Wait()

	// Store should not have panicked and should be in valid state
	notifications := store.List()
	if len(notifications) > 50 {
		t.Errorf("Concurrent operations resulted in more than maxItems: %d", len(notifications))
	}
}

// failingReader always returns an error.
type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errors.New("crypto/rand unavailable")
}

func TestGenerateID_ReturnsValidIDOnSuccess(t *testing.T) {
	t.Parallel()

	id := generateID(rand.Reader)

	if len(id) != 16 {
		t.Errorf("generateID() length = %d, want 16 hex chars", len(id))
	}

	// Verify it's valid hex
	if _, err := hex.DecodeString(id); err != nil {
		t.Errorf("generateID() returned invalid hex: %v", err)
	}
}

func TestGenerateID_ReturnsFallbackIDOnReaderFailure(t *testing.T) {
	t.Parallel()

	// Should NOT panic, should return fallback ID
	id := generateID(failingReader{})

	if id == "" {
		t.Error("generateID() returned empty string on reader failure")
	}

	if len(id) != 16 {
		t.Errorf("generateID() fallback length = %d, want 16 hex chars", len(id))
	}

	// Verify it's valid hex
	if _, err := hex.DecodeString(id); err != nil {
		t.Errorf("generateID() fallback returned invalid hex: %v", err)
	}
}

func TestGenerateID_FallbackIDsAreUnique(t *testing.T) {
	t.Parallel()

	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateID(failingReader{})
		if ids[id] {
			t.Errorf("generateID() produced duplicate fallback ID: %s", id)
		}
		ids[id] = true
	}
}
