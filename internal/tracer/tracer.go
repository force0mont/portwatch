// Package tracer tracks how long each port has been continuously observed
// across scans, emitting a duration for each active port entry.
package tracer

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds the first-seen timestamp and the last-seen timestamp for a port.
type Entry struct {
	FirstSeen time.Time
	LastSeen  time.Time
	Duration  time.Duration
}

// Tracer maintains per-port observation windows.
type Tracer struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a Tracer using the real wall clock.
func New() *Tracer {
	return newWithClock(time.Now)
}

func newWithClock(now func() time.Time) *Tracer {
	return &Tracer{
		entries: make(map[string]Entry),
		now:     now,
	}
}

func key(p scanner.Port) string {
	return p.Protocol + "/" + p.Address + ":" + itoa(p.Port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

// Observe records a port as seen at the current time. Returns the updated Entry.
func (t *Tracer) Observe(p scanner.Port) Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	k := key(p)
	e, ok := t.entries[k]
	if !ok {
		e = Entry{FirstSeen: now}
	}
	e.LastSeen = now
	e.Duration = now.Sub(e.FirstSeen)
	t.entries[k] = e
	return e
}

// Remove deletes the tracking entry for a port (e.g. when it disappears).
func (t *Tracer) Remove(p scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key(p))
}

// Get returns the current entry for a port and whether it exists.
func (t *Tracer) Get(p scanner.Port) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key(p)]
	return e, ok
}

// Len returns the number of actively tracked ports.
func (t *Tracer) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}
