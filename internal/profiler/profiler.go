// Package profiler tracks per-port timing statistics, recording how long
// a port has been continuously open and computing a simple risk score
// based on duration thresholds.
package profiler

import (
	"sync"
	"time"

	"github.com/example/portwatch/internal/scanner"
)

// Entry holds the first-seen timestamp and derived stats for a port.
type Entry struct {
	FirstSeen time.Time
	LastSeen  time.Time
	Duration  time.Duration
	// AgeScore is 0-100; higher means the port has been open longer.
	AgeScore int
}

// Profiler maintains open-duration records for observed ports.
type Profiler struct {
	mu      sync.Mutex
	entries map[string]Entry
	clock   func() time.Time

	// Thresholds used to compute AgeScore.
	warnAfter     time.Duration
	criticalAfter time.Duration
}

// New returns a Profiler with sensible defaults (warn >5 min, critical >1 h).
func New() *Profiler {
	return newWithClock(time.Now, 5*time.Minute, time.Hour)
}

func newWithClock(clock func() time.Time, warn, critical time.Duration) *Profiler {
	return &Profiler{
		entries:       make(map[string]Entry),
		clock:         clock,
		warnAfter:     warn,
		criticalAfter: critical,
	}
}

func key(p scanner.Port) string {
	return p.Protocol + ":" + p.Address + ":" + itoa(p.Port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 6)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

// Observe records a port as seen at the current clock time and returns its Entry.
func (p *Profiler) Observe(port scanner.Port) Entry {
	now := p.clock()
	k := key(port)

	p.mu.Lock()
	defer p.mu.Unlock()

	e, ok := p.entries[k]
	if !ok {
		e = Entry{FirstSeen: now}
	}
	e.LastSeen = now
	e.Duration = now.Sub(e.FirstSeen)
	e.AgeScore = p.score(e.Duration)
	p.entries[k] = e
	return e
}

// Remove deletes the tracking entry for a port (e.g. when it disappears).
func (p *Profiler) Remove(port scanner.Port) {
	p.mu.Lock()
	delete(p.entries, key(port))
	p.mu.Unlock()
}

// Snapshot returns a copy of all current entries.
func (p *Profiler) Snapshot() map[string]Entry {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make(map[string]Entry, len(p.entries))
	for k, v := range p.entries {
		out[k] = v
	}
	return out
}

func (p *Profiler) score(d time.Duration) int {
	switch {
	case d >= p.criticalAfter:
		return 100
	case d >= p.warnAfter:
		return 50
	default:
		return 0
	}
}
