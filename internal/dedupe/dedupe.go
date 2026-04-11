// Package dedupe provides a short-lived deduplication window that suppresses
// identical port events fired in rapid succession.
package dedupe

import (
	"sync"
	"time"

	"github.com/iamcalledned/portwatch/internal/scanner"
)

// Clock is a narrow interface so tests can inject a fake time source.
type Clock func() time.Time

// Dedupe suppresses repeated events for the same (port, protocol, kind)
// tuple within a configurable window.
type Dedupe struct {
	mu     sync.Mutex
	window time.Duration
	clock  Clock
	seen   map[string]time.Time
}

// New returns a Dedupe with the given suppression window.
func New(window time.Duration) *Dedupe {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock Clock) *Dedupe {
	return &Dedupe{
		window: window,
		clock:  clock,
		seen:   make(map[string]time.Time),
	}
}

// IsDuplicate reports whether an equivalent event was already seen within
// the deduplication window. If not, it records the event and returns false.
func (d *Dedupe) IsDuplicate(p scanner.Port, kind string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	d.evict(now)

	key := buildKey(p, kind)
	if _, exists := d.seen[key]; exists {
		return true
	}
	d.seen[key] = now
	return false
}

// Reset clears all recorded events.
func (d *Dedupe) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}

// evict removes entries older than the window. Must be called with mu held.
func (d *Dedupe) evict(now time.Time) {
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

func buildKey(p scanner.Port, kind string) string {
	return p.Protocol + ":" + p.Address + ":" + itoa(p.Port) + ":" + kind
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 6)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
