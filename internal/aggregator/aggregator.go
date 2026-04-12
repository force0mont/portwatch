// Package aggregator batches port events over a fixed time window and
// emits a single combined summary, reducing downstream noise when many
// ports change state simultaneously.
package aggregator

import (
	"sync"
	"time"

	"github.com/joemiller/portwatch/internal/state"
)

// Summary holds all events collected within one flush window.
type Summary struct {
	Appeared []state.PortEvent
	Disappeared []state.PortEvent
	CollectedAt time.Time
}

// Aggregator buffers PortEvents and flushes them as a Summary on demand
// or when the window elapses.
type Aggregator struct {
	mu       sync.Mutex
	buf      []state.PortEvent
	window   time.Duration
	clock    func() time.Time
	lastFlush time.Time
}

// New returns an Aggregator that groups events within the given window.
func New(window time.Duration) *Aggregator {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock func() time.Time) *Aggregator {
	return &Aggregator{
		window:    window,
		clock:     clock,
		lastFlush: clock(),
	}
}

// Add appends a PortEvent to the current buffer.
func (a *Aggregator) Add(e state.PortEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buf = append(a.buf, e)
}

// Ready reports whether the flush window has elapsed since the last flush.
func (a *Aggregator) Ready() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.clock().Sub(a.lastFlush) >= a.window
}

// Flush drains the buffer and returns a Summary. The internal last-flush
// timestamp is reset regardless of whether any events were buffered.
func (a *Aggregator) Flush() Summary {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := a.clock()
	s := Summary{CollectedAt: now}
	for _, e := range a.buf {
		switch e.Kind {
		case state.EventAppeared:
			s.Appeared = append(s.Appeared, e)
		case state.EventDisappeared:
			s.Disappeared = append(s.Disappeared, e)
		}
	}
	a.buf = a.buf[:0]
	a.lastFlush = now
	return s
}
