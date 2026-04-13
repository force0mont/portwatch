// Package window provides a sliding-window counter keyed by an arbitrary
// string. It is safe for concurrent use.
package window

import (
	"sync"
	"time"
)

// Clock allows the wall-clock to be replaced in tests.
type Clock func() time.Time

type bucket struct {
	count int
	at    time.Time
}

// Window tracks how many times a key has been seen within a rolling duration.
type Window struct {
	mu       sync.Mutex
	size     time.Duration
	clock    Clock
	buckets  map[string][]bucket
}

// New returns a Window with the given rolling duration.
func New(size time.Duration) *Window {
	return newWithClock(size, time.Now)
}

func newWithClock(size time.Duration, clock Clock) *Window {
	return &Window{
		size:    size,
		clock:   clock,
		buckets: make(map[string][]bucket),
	}
}

// Add records one occurrence for key and returns the total count within the
// current window.
func (w *Window) Add(key string) int {
	now := w.clock()
	w.mu.Lock()
	defer w.mu.Unlock()

	w.evict(key, now)
	w.buckets[key] = append(w.buckets[key], bucket{count: 1, at: now})
	return w.sum(key)
}

// Count returns the number of occurrences for key within the current window
// without modifying state.
func (w *Window) Count(key string) int {
	now := w.clock()
	w.mu.Lock()
	defer w.mu.Unlock()

	w.evict(key, now)
	return w.sum(key)
}

// Reset clears all recorded buckets for key.
func (w *Window) Reset(key string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.buckets, key)
}

// evict removes buckets that have fallen outside the window. Must be called
// with w.mu held.
func (w *Window) evict(key string, now time.Time) {
	cutoff := now.Add(-w.size)
	bs := w.buckets[key]
	i := 0
	for i < len(bs) && bs[i].at.Before(cutoff) {
		i++
	}
	w.buckets[key] = bs[i:]
}

// sum returns the total count across all retained buckets. Must be called
// with w.mu held.
func (w *Window) sum(key string) int {
	total := 0
	for _, b := range w.buckets[key] {
		total += b.count
	}
	return total
}
