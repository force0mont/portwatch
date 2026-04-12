// Package trimmer evicts the oldest entries from a bounded port-event
// buffer, keeping memory usage predictable over long daemon runs.
package trimmer

import (
	"sync"
	"time"
)

// Entry is a timestamped item held in the buffer.
type Entry struct {
	At    time.Time
	Key   string
	Value any
}

// Trimmer keeps at most Cap entries, evicting the oldest when the cap is
// exceeded or when entries older than TTL are pruned.
type Trimmer struct {
	mu      sync.Mutex
	entries []Entry
	cap     int
	ttl     time.Duration
	now     func() time.Time
}

// New returns a Trimmer with the given capacity and TTL.
func New(cap int, ttl time.Duration) *Trimmer {
	return newWithClock(cap, ttl, time.Now)
}

func newWithClock(cap int, ttl time.Duration, now func() time.Time) *Trimmer {
	if cap <= 0 {
		cap = 256
	}
	return &Trimmer{cap: cap, ttl: ttl, now: now}
}

// Add appends an entry then enforces both the TTL and the hard cap.
func (t *Trimmer) Add(key string, value any) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.entries = append(t.entries, Entry{At: t.now(), Key: key, Value: value})
	t.evictLocked()
}

// Len returns the current number of buffered entries.
func (t *Trimmer) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}

// All returns a snapshot copy of all live entries.
func (t *Trimmer) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, len(t.entries))
	copy(out, t.entries)
	return out
}

// Prune removes entries older than TTL without adding a new one.
func (t *Trimmer) Prune() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evictLocked()
}

// evictLocked must be called with t.mu held.
func (t *Trimmer) evictLocked() {
	// Drop expired entries first.
	if t.ttl > 0 {
		cutoff := t.now().Add(-t.ttl)
		i := 0
		for i < len(t.entries) && t.entries[i].At.Before(cutoff) {
			i++
		}
		t.entries = t.entries[i:]
	}
	// Enforce hard cap by dropping oldest.
	if len(t.entries) > t.cap {
		t.entries = t.entries[len(t.entries)-t.cap:]
	}
}
