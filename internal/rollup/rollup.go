// Package rollup groups repeated alert events into a single summary
// event after a configurable count threshold is reached within a window.
package rollup

import (
	"fmt"
	"sync"
	"time"
)

// Entry tracks occurrences of a unique event key within a window.
type Entry struct {
	Count     int
	WindowEnd time.Time
	Flushed   bool
}

// Rollup accumulates events and emits a summary once the threshold is met.
type Rollup struct {
	mu        sync.Mutex
	clock     func() time.Time
	window    time.Duration
	threshold int
	entries   map[string]*Entry
}

// New creates a Rollup with the given window duration and count threshold.
func New(window time.Duration, threshold int) *Rollup {
	return newWithClock(window, threshold, time.Now)
}

func newWithClock(window time.Duration, threshold int, clock func() time.Time) *Rollup {
	return &Rollup{
		clock:     clock,
		window:    window,
		threshold: threshold,
		entries:   make(map[string]*Entry),
	}
}

// Record records an occurrence of key. It returns (summary, true) when the
// threshold is first exceeded within the current window, otherwise ("", false).
func (r *Rollup) Record(key string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.clock()
	e, ok := r.entries[key]
	if !ok || now.After(e.WindowEnd) {
		r.entries[key] = &Entry{
			Count:     1,
			WindowEnd: now.Add(r.window),
			Flushed:   false,
		}
		return "", false
	}

	e.Count++
	if e.Count == r.threshold && !e.Flushed {
		e.Flushed = true
		return fmt.Sprintf("%s repeated %d times in %s", key, e.Count, r.window), true
	}
	return "", false
}

// Reset clears all tracked state.
func (r *Rollup) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = make(map[string]*Entry)
}
