// Package decay implements a score-decay tracker that reduces the
// accumulated risk score of a port over time when no new alerts are seen.
package decay

import (
	"math"
	"sync"
	"time"
)

// clock allows deterministic testing.
type clock func() time.Time

// entry holds the current score and the last time it was updated.
type entry struct {
	score     float64
	updatedAt time.Time
}

// Tracker reduces a port's accumulated score toward zero at a configurable
// half-life. Each call to Add bumps the score; Score returns the current
// decayed value.
type Tracker struct {
	mu       sync.Mutex
	entries  map[string]*entry
	halfLife time.Duration
	now      clock
}

// New returns a Tracker with the given half-life duration.
func New(halfLife time.Duration) *Tracker {
	return newWithClock(halfLife, time.Now)
}

func newWithClock(halfLife time.Duration, c clock) *Tracker {
	return &Tracker{
		entries:  make(map[string]*entry),
		halfLife: halfLife,
		now:      c,
	}
}

// Add increases the score for key by delta, after first applying decay since
// the last update.
func (t *Tracker) Add(key string, delta float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	e, ok := t.entries[key]
	if !ok {
		t.entries[key] = &entry{score: delta, updatedAt: now}
		return
	}
	e.score = t.decayed(e, now) + delta
	e.updatedAt = now
}

// Score returns the current decayed score for key, or 0 if unknown.
func (t *Tracker) Score(key string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.entries[key]
	if !ok {
		return 0
	}
	return t.decayed(e, t.now())
}

// Reset removes the tracked score for key.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// decayed computes the exponentially-decayed score at time now.
func (t *Tracker) decayed(e *entry, now time.Time) float64 {
	elapsed := now.Sub(e.updatedAt)
	if elapsed <= 0 || t.halfLife <= 0 {
		return e.score
	}
	exponent := float64(elapsed) / float64(t.halfLife)
	return e.score * math.Pow(0.5, exponent)
}
