// Package budget implements a per-key token-bucket style alert budget.
// A Budget caps the total number of alerts that may be emitted for a
// given key within a rolling time window.  Once the cap is reached,
// further calls to Allow return false until the window expires.
package budget

import (
	"sync"
	"time"
)

// clock is a seam for deterministic tests.
type clock func() time.Time

// entry tracks usage for a single key.
type entry struct {
	count     int
	windowEnd time.Time
}

// Budget is a concurrency-safe per-key alert budget.
type Budget struct {
	mu     sync.Mutex
	clock  clock
	window time.Duration
	max    int
	keys   map[string]*entry
}

// New returns a Budget that allows at most max events per key within
// the given window duration.
func New(max int, window time.Duration) *Budget {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, c clock) *Budget {
	if max <= 0 {
		max = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &Budget{
		clock:  c,
		window: window,
		max:    max,
		keys:   make(map[string]*entry),
	}
}

// Allow returns true and increments the counter if the key is within
// budget.  It returns false once the cap has been reached for the
// current window.
func (b *Budget) Allow(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.clock()
	e, ok := b.keys[key]
	if !ok || now.After(e.windowEnd) {
		b.keys[key] = &entry{count: 1, windowEnd: now.Add(b.window)}
		return true
	}
	if e.count >= b.max {
		return false
	}
	e.count++
	return true
}

// Remaining returns how many more events are allowed for key in the
// current window.  A return value of 0 means the budget is exhausted.
func (b *Budget) Remaining(key string) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.clock()
	e, ok := b.keys[key]
	if !ok || now.After(e.windowEnd) {
		return b.max
	}
	r := b.max - e.count
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears the budget for key, allowing a fresh window to begin on
// the next call to Allow.
func (b *Budget) Reset(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.keys, key)
}
