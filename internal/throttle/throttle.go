// Package throttle provides per-key event throttling to prevent
// notification storms when a port flaps repeatedly.
package throttle

import (
	"fmt"
	"sync"
	"time"
)

// clock abstracts time for testing.
type clock func() time.Time

// entry tracks the count and window start for a single key.
type entry struct {
	count     int
	windowEnd time.Time
}

// Throttle limits how many events are forwarded per key within a window.
type Throttle struct {
	mu       sync.Mutex
	entries  map[string]*entry
	max      int
	window   time.Duration
	nowFn    clock
}

// New returns a Throttle that allows at most max events per key per window.
func New(max int, window time.Duration) *Throttle {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, fn clock) *Throttle {
	return &Throttle{
		entries: make(map[string]*entry),
		max:     max,
		window:  window,
		nowFn:   fn,
	}
}

// Allow returns true if the event for the given port+protocol should be
// forwarded, false if it has been throttled.
func (t *Throttle) Allow(port uint16, proto string) bool {
	key := fmt.Sprintf("%s:%d", proto, port)
	now := t.nowFn()

	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.entries[key]
	if !ok || now.After(e.windowEnd) {
		t.entries[key] = &entry{count: 1, windowEnd: now.Add(t.window)}
		return true
	}

	e.count++
	return e.count <= t.max
}

// Reset clears all throttle state (useful for testing or config reload).
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*entry)
}
