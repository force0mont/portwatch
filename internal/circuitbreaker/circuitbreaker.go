// Package circuitbreaker provides a simple circuit breaker that opens after
// a configurable number of consecutive failures, preventing further calls
// until a cooldown period has elapsed.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit is open and calls are rejected.
var ErrOpen = errors.New("circuit breaker is open")

// state represents the current circuit breaker state.
type state int

const (
	stateClosed state = iota
	stateOpen
)

// clock allows time injection for testing.
type clock func() time.Time

// Breaker is a circuit breaker that opens after maxFailures consecutive
// failures and resets after the cooldown window.
type Breaker struct {
	mu          sync.Mutex
	state       state
	failures    int
	maxFailures int
	cooldown    time.Duration
	openedAt    time.Time
	now         clock
}

// New returns a Breaker that opens after maxFailures consecutive failures
// and attempts recovery after cooldown.
func New(maxFailures int, cooldown time.Duration) *Breaker {
	return newWithClock(maxFailures, cooldown, time.Now)
}

func newWithClock(maxFailures int, cooldown time.Duration, now clock) *Breaker {
	return &Breaker{
		maxFailures: maxFailures,
		cooldown:    cooldown,
		now:         now,
	}
}

// Allow returns nil if the call is permitted, or ErrOpen if the circuit is
// open and the cooldown has not yet elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == stateOpen {
		if b.now().Sub(b.openedAt) >= b.cooldown {
			// half-open: allow one probe
			b.state = stateClosed
			b.failures = 0
			return nil
		}
		return ErrOpen
	}
	return nil
}

// RecordSuccess resets the failure counter.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = stateClosed
}

// RecordFailure increments the failure counter and opens the circuit if the
// threshold has been reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = stateOpen
		b.openedAt = b.now()
	}
}

// IsOpen reports whether the circuit is currently open.
func (b *Breaker) IsOpen() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state == stateOpen
}
