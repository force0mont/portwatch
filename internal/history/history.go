// Package history maintains a rolling window of port-scan events so that
// callers can query recent activity without re-scanning.
package history

import (
	"sync"
	"time"

	"github.com/yourorg/portwatch/internal/state"
)

// Entry pairs a state event with the time it was recorded.
type Entry struct {
	RecordedAt time.Time
	Event      state.Event
}

// History stores the most recent scan events up to a configurable capacity.
type History struct {
	mu       sync.RWMutex
	entries  []Entry
	capacity int
}

// New returns a History that retains at most capacity entries.
// If capacity is <= 0 it defaults to 100.
func New(capacity int) *History {
	if capacity <= 0 {
		capacity = 100
	}
	return &History{
		entries:  make([]Entry, 0, capacity),
		capacity: capacity,
	}
}

// Record appends a batch of events to the history, evicting the oldest
// entries when the capacity would be exceeded.
func (h *History) Record(events []state.Event) {
	if len(events) == 0 {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now().UTC()
	for _, e := range events {
		if len(h.entries) == h.capacity {
			// shift left – O(n) but capacity is small
			h.entries = h.entries[1:]
		}
		h.entries = append(h.entries, Entry{RecordedAt: now, Event: e})
	}
}

// Since returns all entries recorded at or after t.
func (h *History) Since(t time.Time) []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var out []Entry
	for _, e := range h.entries {
		if !e.RecordedAt.Before(t) {
			out = append(out, e)
		}
	}
	return out
}

// All returns a snapshot of every entry currently held.
func (h *History) All() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Len returns the number of entries currently stored.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}
