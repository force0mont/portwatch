package escalator

import (
	"sync"
	"time"
)

// Level represents an escalation severity tier.
type Level int

const (
	LevelNone    Level = iota
	LevelWarning       // repeated within short window
	LevelCritical      // repeated beyond threshold
)

// Entry tracks how many times a key has triggered within a window.
type Entry struct {
	count     int
	windowEnd time.Time
}

// Escalator promotes alert severity based on repeated occurrences within a
// rolling time window. Keys are arbitrary strings (e.g. "tcp:8080").
type Escalator struct {
	mu             sync.Mutex
	entries        map[string]*Entry
	window         time.Duration
	warningAt      int
	criticalAt     int
	clock          func() time.Time
}

// New returns an Escalator with the given window and thresholds.
// warningAt is the hit count that triggers LevelWarning;
// criticalAt triggers LevelCritical.
func New(window time.Duration, warningAt, criticalAt int) *Escalator {
	return newWithClock(window, warningAt, criticalAt, time.Now)
}

func newWithClock(window time.Duration, warningAt, criticalAt int, clock func() time.Time) *Escalator {
	return &Escalator{
		entries:    make(map[string]*Entry),
		window:     window,
		warningAt:  warningAt,
		criticalAt: criticalAt,
		clock:      clock,
	}
}

// Record increments the hit count for key and returns the resulting Level.
func (e *Escalator) Record(key string) Level {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.clock()
	ent, ok := e.entries[key]
	if !ok || now.After(ent.windowEnd) {
		ent = &Entry{windowEnd: now.Add(e.window)}
		e.entries[key] = ent
	}
	ent.count++

	switch {
	case ent.count >= e.criticalAt:
		return LevelCritical
	case ent.count >= e.warningAt:
		return LevelWarning
	default:
		return LevelNone
	}
}

// Reset clears the state for key.
func (e *Escalator) Reset(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.entries, key)
}

// Len returns the number of tracked keys.
func (e *Escalator) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.entries)
}
