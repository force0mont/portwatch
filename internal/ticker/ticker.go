// Package ticker provides a configurable scan interval ticker that supports
// jittered scheduling to avoid thundering-herd effects when multiple instances
// run concurrently.
package ticker

import (
	"context"
	"time"
)

// Ticker fires at a configurable interval, optionally applying jitter.
type Ticker struct {
	interval time.Duration
	jitter   func(base time.Duration) time.Duration
	now      func() time.Time
}

// Option configures a Ticker.
type Option func(*Ticker)

// WithJitter sets a jitter function applied to each interval before sleeping.
func WithJitter(fn func(base time.Duration) time.Duration) Option {
	return func(t *Ticker) { t.jitter = fn }
}

// withNow overrides the clock (for testing).
func withNow(fn func() time.Time) Option {
	return func(t *Ticker) { t.now = fn }
}

// New creates a Ticker that fires every interval.
// interval must be > 0, otherwise New panics.
func New(interval time.Duration, opts ...Option) *Ticker {
	if interval <= 0 {
		panic("ticker: interval must be positive")
	}
	t := &Ticker{
		interval: interval,
		jitter:   func(d time.Duration) time.Duration { return d },
		now:      time.Now,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Run calls fn each time the ticker fires until ctx is cancelled.
// The first call to fn happens after one interval (not immediately).
func (t *Ticker) Run(ctx context.Context, fn func(at time.Time)) error {
	for {
		delay := t.jitter(t.interval)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			fn(t.now())
		}
	}
}
