// Package observer tracks how often each port is seen across scans
// and emits a stability score indicating how consistently a port appears.
package observer

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds observation data for a single port.
type Entry struct {
	Port      scanner.Port
	FirstSeen time.Time
	LastSeen  time.Time
	SeenCount int
	MissCount int
}

// StabilityScore returns a value in [0.0, 1.0] representing how reliably
// the port has been present across all scans observed so far.
func (e Entry) StabilityScore() float64 {
	total := e.SeenCount + e.MissCount
	if total == 0 {
		return 0
	}
	return float64(e.SeenCount) / float64(total)
}

type clock func() time.Time

// Observer records presence or absence of ports across scan cycles.
type Observer struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	now     clock
}

// New returns an Observer using the real wall clock.
func New() *Observer {
	return newWithClock(time.Now)
}

func newWithClock(c clock) *Observer {
	return &Observer{
		entries: make(map[string]*Entry),
		now:     c,
	}
}

func key(p scanner.Port) string {
	return p.Protocol + ":" + p.Address
}

// Observe records which ports were seen in the latest scan.
// Ports present in seen are marked as present; all previously known ports
// absent from seen have their miss counter incremented.
func (o *Observer) Observe(seen []scanner.Port) {
	o.mu.Lock()
	defer o.mu.Unlock()

	now := o.now()
	seen	Map := make(map[string]scanner.Port, len(seen))
	for _, p := range seen {
		seenMap[key(p)] = p
	}

	for k, e := range o.entries {
		if _, ok := seenMap[k]; !ok {
			e.MissCount++
		}
	}

	for k, p := range seenMap {
		if e, ok := o.entries[k]; ok {
			e.SeenCount++
			e.LastSeen = now
		} else {
			o.entries[k] = &Entry{
				Port:      p,
				FirstSeen: now,
				LastSeen:  now,
				SeenCount: 1,
			}
		}
	}
}

// Get returns the Entry for the given port and whether it was found.
func (o *Observer) Get(p scanner.Port) (Entry, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	e, ok := o.entries[key(p)]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot copy of all tracked entries.
func (o *Observer) All() []Entry {
	o.mu.RLock()
	defer o.mu.RUnlock()
	out := make([]Entry, 0, len(o.entries))
	for _, e := range o.entries {
		out = append(out, *e)
	}
	return out
}

// Len returns the number of tracked ports.
func (o *Observer) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.entries)
}
