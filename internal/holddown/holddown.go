// Package holddown suppresses repeated alerts for a port until it has been
// absent for a minimum quiet period, preventing flapping notifications.
package holddown

import (
	"fmt"
	"sync"
	"time"
)

// Clock is a time source used for testing.
type Clock func() time.Time

// entry tracks the last time a port was seen and whether it is held down.
type entry struct {
	lastSeen  time.Time
	held      bool
}

// HoldDown suppresses re-alerts until a port has been quiet for QuietFor.
type HoldDown struct {
	mu      sync.Mutex
	entries map[string]*entry
	quietFor time.Duration
	now     Clock
}

// New creates a HoldDown with the given quiet period.
func New(quietFor time.Duration) *HoldDown {
	return newWithClock(quietFor, time.Now)
}

func newWithClock(quietFor time.Duration, now Clock) *HoldDown {
	return &HoldDown{
		entries:  make(map[string]*entry),
		quietFor: quietFor,
		now:      now,
	}
}

func key(port uint16, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Seen records that the port is currently active. Returns true if this is the
// first time the port is seen (i.e. alert should fire).
func (h *HoldDown) Seen(port uint16, proto string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	k := key(port, proto)
	now := h.now()
	e, exists := h.entries[k]
	if !exists {
		h.entries[k] = &entry{lastSeen: now, held: true}
		return true
	}
	e.lastSeen = now
	return false
}

// Gone records that the port is no longer visible. Once the port has been
// gone for the quiet period, Seen will fire again on re-appearance.
func (h *HoldDown) Gone(port uint16, proto string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.entries, key(port, proto))
}

// Prune removes entries whose lastSeen is older than quietFor, releasing the
// hold so a re-appearance will trigger a fresh alert.
func (h *HoldDown) Prune() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	cutoff := h.now().Add(-h.quietFor)
	removed := 0
	for k, e := range h.entries {
		if e.lastSeen.Before(cutoff) {
			delete(h.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of currently tracked ports.
func (h *HoldDown) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.entries)
}
