// Package planner schedules the next scan time based on interval, jitter,
// and backpressure signals from the watcher pipeline.
package planner

import (
	"sync"
	"time"
)

// Clock is a minimal time source used for testing.
type Clock func() time.Time

// Planner computes the next scan deadline and tracks whether the previous
// scan completed within its allotted window.
type Planner struct {
	mu       sync.Mutex
	now      Clock
	interval time.Duration
	jitter   time.Duration
	lastFire time.Time
	missed   int
}

// New returns a Planner with the given base interval and maximum jitter.
// jitter must be less than interval; if it is not, it is clamped to zero.
func New(interval, jitter time.Duration) *Planner {
	if jitter >= interval {
		jitter = 0
	}
	return newWithClock(interval, jitter, time.Now)
}

func newWithClock(interval, jitter time.Duration, now Clock) *Planner {
	return &Planner{
		now:      now,
		interval: interval,
		jitter:   jitter,
	}
}

// Next returns the duration to wait before the next scan fires.
// It applies a pseudo-deterministic jitter derived from the current
// nanosecond timestamp.
func (p *Planner) Next() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()

	base := p.interval
	if p.jitter > 0 {
		offset := time.Duration(p.now().UnixNano()%int64(p.jitter)) - p.jitter/2
		base += offset
		if base <= 0 {
			base = p.interval
		}
	}
	return base
}

// Mark records that a scan fired at time t. It updates the missed-scan
// counter if the gap since the last fire exceeds 1.5× the interval.
func (p *Planner) Mark(t time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.lastFire.IsZero() {
		gap := t.Sub(p.lastFire)
		if gap > p.interval+p.interval/2 {
			p.missed++
		}
	}
	p.lastFire = t
}

// Missed returns the number of scans that were detected as late.
func (p *Planner) Missed() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.missed
}

// Reset clears the missed counter and last-fire timestamp.
func (p *Planner) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.missed = 0
	p.lastFire = time.Time{}
}
