// Package fencepost tracks the last-seen timestamp for a named checkpoint,
// allowing the daemon to detect gaps in scan coverage (e.g. missed intervals).
package fencepost

import (
	"sync"
	"time"
)

// Clock is a func that returns the current time (injectable for tests).
type Clock func() time.Time

// Post records the last-marked time for a named checkpoint and exposes
// whether the gap since the previous mark exceeds a configurable threshold.
type Post struct {
	mu        sync.Mutex
	clock     Clock
	marks     map[string]time.Time
	threshold time.Duration
}

// New returns a Post with the given missed-interval threshold.
func New(threshold time.Duration) *Post {
	return newWithClock(threshold, time.Now)
}

func newWithClock(threshold time.Duration, clock Clock) *Post {
	return &Post{
		clock:     clock,
		marks:     make(map[string]time.Time),
		threshold: threshold,
	}
}

// Mark records the current time for the named checkpoint.
func (p *Post) Mark(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.marks[name] = p.clock()
}

// Gap returns the duration since the last mark for name.
// If name has never been marked, Gap returns 0 and ok == false.
func (p *Post) Gap(name string) (time.Duration, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	t, ok := p.marks[name]
	if !ok {
		return 0, false
	}
	return p.clock().Sub(t), true
}

// Overdue returns true when the gap for name exceeds the configured threshold.
// An unseen name is never considered overdue.
func (p *Post) Overdue(name string) bool {
	gap, ok := p.Gap(name)
	return ok && gap > p.threshold
}

// Reset removes the checkpoint for name.
func (p *Post) Reset(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.marks, name)
}
