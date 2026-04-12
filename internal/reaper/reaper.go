// Package reaper removes stale port entries that have not been seen
// within a configurable TTL, preventing unbounded memory growth in
// long-running portwatch daemons.
package reaper

import (
	"sync"
	"time"
)

// clock is a seam for deterministic testing.
type clock func() time.Time

// Entry tracks the last time a port key was observed.
type Entry struct {
	LastSeen time.Time
	Count    int
}

// Reaper evicts entries that have not been refreshed within TTL.
type Reaper struct {
	mu      sync.Mutex
	entries map[string]Entry
	ttl     time.Duration
	now     clock
}

// New returns a Reaper with the given TTL.
func New(ttl time.Duration) *Reaper {
	return newWithClock(ttl, time.Now)
}

func newWithClock(ttl time.Duration, now clock) *Reaper {
	return &Reaper{
		entries: make(map[string]Entry),
		ttl:     ttl,
		now:     now,
	}
}

// Touch records or refreshes an entry for key.
func (r *Reaper) Touch(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e := r.entries[key]
	e.LastSeen = r.now()
	e.Count++
	r.entries[key] = e
}

// Reap removes all entries older than TTL and returns their keys.
func (r *Reaper) Reap() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := r.now().Add(-r.ttl)
	var evicted []string
	for k, e := range r.entries {
		if e.LastSeen.Before(cutoff) {
			evicted = append(evicted, k)
			delete(r.entries, k)
		}
	}
	return evicted
}

// Len returns the number of tracked entries.
func (r *Reaper) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}
