// Package limiter provides a token-bucket style per-key rate limiter that
// caps the number of alert events emitted for a given port+protocol pair
// within a sliding time window.
package limiter

import (
	"fmt"
	"sync"
	"time"
)

// clock allows deterministic testing.
type clock func() time.Time

// bucket tracks token state for a single key.
type bucket struct {
	tokens    int
	windowEnd time.Time
}

// Limiter enforces a maximum number of events per key per window.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	max     int
	window  time.Duration
	now     clock
}

// New returns a Limiter that allows at most max events per window duration
// for each unique port+protocol key.
func New(max int, window time.Duration) *Limiter {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, now clock) *Limiter {
	return &Limiter{
		buckets: make(map[string]*bucket),
		max:     max,
		window:  window,
		now:     now,
	}
}

// Allow returns true if the event for the given port and protocol should be
// allowed through. Once max events have been emitted in the current window
// the key is suppressed until the window rolls over.
func (l *Limiter) Allow(port uint16, proto string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := fmt.Sprintf("%s:%d", proto, port)
	now := l.now()

	b, ok := l.buckets[key]
	if !ok || now.After(b.windowEnd) {
		l.buckets[key] = &bucket{tokens: 1, windowEnd: now.Add(l.window)}
		return true
	}

	if b.tokens >= l.max {
		return false
	}

	b.tokens++
	return true
}

// Reset clears all bucket state, useful for testing or forced resets.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = make(map[string]*bucket)
}
