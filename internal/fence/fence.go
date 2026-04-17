// Package fence provides a simple trip-wire that fires once a counter
// exceeds a configured threshold within a sliding time window.
package fence

import (
	"sync"
	"time"
)

// Tripped is returned by Check when the threshold is crossed.
type Tripped struct {
	Key   string
	Count int
	At    time.Time
}

type entry struct {
	times []time.Time
}

// Fence tracks per-key event counts and signals when a threshold is breached.
type Fence struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	entries   map[string]*entry
	now       func() time.Time
}

// New returns a Fence with the given sliding window and threshold.
func New(window time.Duration, threshold int) *Fence {
	return newWithClock(window, threshold, time.Now)
}

func newWithClock(window time.Duration, threshold int, now func() time.Time) *Fence {
	return &Fence{
		window:    window,
		threshold: threshold,
		entries:   make(map[string]*entry),
		now:       now,
	}
}

// Record adds an event for key and returns a Tripped value if the threshold
// is reached, along with true. Returns zero Tripped and false otherwise.
func (f *Fence) Record(key string) (Tripped, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := f.now()
	cutoff := now.Add(-f.window)

	e, ok := f.entries[key]
	if !ok {
		e = &entry{}
		f.entries[key] = e
	}

	// evict old timestamps
	filtered := e.times[:0]
	for _, t := range e.times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	e.times = append(filtered, now)

	if len(e.times) >= f.threshold {
		return Tripped{Key: key, Count: len(e.times), At: now}, true
	}
	return Tripped{}, false
}

// Reset clears the event history for key.
func (f *Fence) Reset(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.entries, key)
}
