// Package debounce provides a port-event debouncer that suppresses
// transient appearances and disappearances of ports within a configurable
// stabilisation window. Only events that persist for the full window are
// forwarded to the caller.
package debounce

import (
	"sync"
	"time"

	"github.com/rjbrown57/portwatch/internal/state"
)

// Clock is a thin interface so tests can inject a fake time source.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// entry tracks when a port event was first observed.
type entry struct {
	event     state.Event
	firstSeen time.Time
}

// Debouncer holds pending events and releases them only after they have
// been seen continuously for at least Window duration.
type Debouncer struct {
	mu      sync.Mutex
	Window  time.Duration
	clock   Clock
	pending map[string]entry // key → entry
}

// New returns a Debouncer with the given stabilisation window.
func New(window time.Duration) *Debouncer {
	return newWithClock(window, realClock{})
}

func newWithClock(window time.Duration, c Clock) *Debouncer {
	return &Debouncer{
		Window:  window,
		clock:   c,
		pending: make(map[string]entry),
	}
}

// Feed accepts a batch of new events. It returns the subset that have been
// pending for at least Window and removes them from the pending set.
// Events not yet in the pending set are added; events absent from the new
// batch are evicted (the transient condition has resolved itself).
func (d *Debouncer) Feed(events []state.Event) []state.Event {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock.Now()

	// Build a lookup of the incoming events.
	incoming := make(map[string]state.Event, len(events))
	for _, e := range events {
		k := key(e)
		incoming[k] = e
	}

	// Evict pending entries that are no longer present.
	for k := range d.pending {
		if _, ok := incoming[k]; !ok {
			delete(d.pending, k)
		}
	}

	// Add new arrivals; collect those that have stabilised.
	var stable []state.Event
	for k, e := range incoming {
		if en, exists := d.pending[k]; exists {
			if now.Sub(en.firstSeen) >= d.Window {
				stable = append(stable, e)
				delete(d.pending, k)
			}
		} else {
			d.pending[k] = entry{event: e, firstSeen: now}
		}
	}

	return stable
}

// Len returns the number of events currently held in the pending set.
func (d *Debouncer) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.pending)
}

func key(e state.Event) string {
	return e.Protocol + "|" + e.Addr
}
