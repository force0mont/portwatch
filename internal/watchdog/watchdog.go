package watchdog

import (
	"context"
	"sync"
	"time"
)

// Watchdog monitors the health of the scan loop and emits a signal if
// no heartbeat is received within the configured timeout.
type Watchdog struct {
	timeout  time.Duration
	ticker   func() time.Time
	mu       sync.Mutex
	lastBeat time.Time
	stuck    chan struct{}
}

// New creates a Watchdog that fires if no heartbeat arrives within timeout.
func New(timeout time.Duration) *Watchdog {
	return newWithClock(timeout, time.Now)
}

func newWithClock(timeout time.Duration, clock func() time.Time) *Watchdog {
	return &Watchdog{
		timeout:  timeout,
		ticker:   clock,
		lastBeat: clock(),
		stuck:    make(chan struct{}, 1),
	}
}

// Beat records a heartbeat, resetting the watchdog timer.
func (w *Watchdog) Beat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastBeat = w.ticker()
}

// Stuck returns a channel that receives a signal when the watchdog fires.
func (w *Watchdog) Stuck() <-chan struct{} {
	return w.stuck
}

// Run starts the watchdog loop. It blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context) {
	interval := w.timeout / 2
	if interval < time.Millisecond {
		interval = time.Millisecond
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			w.mu.Lock()
			elapsed.ticker().Sub(w.lastBeat)
			wif elapsed >= w.timeout 			case w.stuck <- struct{}{}:
				default:
				}
			}
		}
	}
}
