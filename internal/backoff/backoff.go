// Package backoff provides an exponential back-off helper used when
// repeated failures occur (e.g. webhook delivery, scanner errors).
// The delay doubles on each consecutive failure up to a configurable
// maximum, and resets to the base delay after a successful call.
package backoff

import (
	"sync"
	"time"
)

// Clock is a narrow interface so tests can inject a fake time source.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// Backoff tracks per-key exponential back-off state.
type Backoff struct {
	mu      sync.Mutex
	base    time.Duration
	max     time.Duration
	clock   Clock
	entries map[string]*entry
}

type entry struct {
	failures int
	next     time.Time
}

// New returns a Backoff with the given base and maximum delay.
func New(base, max time.Duration) *Backoff {
	return newWithClock(base, max, realClock{})
}

func newWithClock(base, max time.Duration, clk Clock) *Backoff {
	return &Backoff{
		base:    base,
		max:     max,
		clock:   clk,
		entries: make(map[string]*entry),
	}
}

// Ready reports whether the key is allowed to proceed right now.
func (b *Backoff) Ready(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.entries[key]
	if !ok {
		return true
	}
	return !b.clock.Now().Before(e.next)
}

// Failure records a failure for key and advances its next-allowed time.
func (b *Backoff) Failure(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.entries[key]
	if !ok {
		e = &entry{}
		b.entries[key] = e
	}
	e.failures++
	delay := b.base
	for i := 1; i < e.failures; i++ {
		delay *= 2
		if delay > b.max {
			delay = b.max
			break
		}
	}
	e.next = b.clock.Now().Add(delay)
}

// Success resets the back-off state for key.
func (b *Backoff) Success(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.entries, key)
}

// Failures returns the current consecutive failure count for key.
func (b *Backoff) Failures(key string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	if e, ok := b.entries[key]; ok {
		return e.failures
	}
	return 0
}
