// Package cooldown provides a per-key cooldown tracker that prevents
// repeated actions from firing more frequently than a configured duration.
package cooldown

import (
	"sync"
	"time"
)

// clock allows injecting a fake time source in tests.
type clock func() time.Time

// Tracker tracks the last time an action was taken for a given key and
// reports whether the cooldown period has elapsed.
type Tracker struct {
	mu       sync.Mutex
	duration time.Duration
	last     map[string]time.Time
	now      clock
}

// New returns a Tracker with the given cooldown duration.
func New(d time.Duration) *Tracker {
	return newWithClock(d, time.Now)
}

func newWithClock(d time.Duration, c clock) *Tracker {
	return &Tracker{
		duration: d,
		last:     make(map[string]time.Time),
		now:      c,
	}
}

// Ready reports whether the cooldown for key has elapsed since the last call
// to Mark. It returns true on the first call for a previously unseen key.
func (t *Tracker) Ready(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	last, seen := t.last[key]
	if !seen {
		return true
	}
	return t.now().Sub(last) >= t.duration
}

// Mark records the current time as the last action time for key.
func (t *Tracker) Mark(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last[key] = t.now()
}

// ReadyAndMark atomically checks whether the key is ready and, if so, marks it.
// Returns true if the action should proceed.
func (t *Tracker) ReadyAndMark(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	last, seen := t.last[key]
	if seen && now.Sub(last) < t.duration {
		return false
	}
	t.last[key] = now
	return true
}

// Reset clears the cooldown state for key, allowing the next call to Ready
// to return true regardless of when Mark was last called.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}
