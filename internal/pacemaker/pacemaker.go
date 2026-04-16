// Package pacemaker tracks scan intervals and flags when scans are
// running slower than the configured target, exposing a simple
// health signal to the rest of the pipeline.
package pacemaker

import (
	"sync"
	"time"
)

// Pacemaker records scan completion times and reports whether the
// observed interval exceeds the allowed threshold.
type Pacemaker struct {
	mu        sync.Mutex
	clock     func() time.Time
	last      time.Time
	threshold time.Duration
	missed    int
}

// New returns a Pacemaker that considers any interval longer than
// threshold to be a missed beat.
func New(threshold time.Duration) *Pacemaker {
	return newWithClock(threshold, time.Now)
}

func newWithClock(threshold time.Duration, clock func() time.Time) *Pacemaker {
	return &Pacemaker{threshold: threshold, clock: clock}
}

// Beat records that a scan completed at the current clock time.
// It returns true if the interval since the previous beat was within
// the allowed threshold (or this is the first beat).
func (p *Pacemaker) Beat() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.clock()
	if p.last.IsZero() {
		p.last = now
		return true
	}
	ok := now.Sub(p.last) <= p.threshold
	if !ok {
		p.missed++
	}
	p.last = now
	return ok
}

// Missed returns the number of intervals that exceeded the threshold.
func (p *Pacemaker) Missed() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.missed
}

// Reset clears the beat history and missed counter.
func (p *Pacemaker) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.last = time.Time{}
	p.missed = 0
}

// LastBeat returns the time of the most recent beat, or the zero
// value if Beat has never been called.
func (p *Pacemaker) LastBeat() time.Time {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.last
}
