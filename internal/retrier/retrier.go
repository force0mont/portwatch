// Package retrier provides a generic retry mechanism with configurable
// attempt limits, delay strategy, and context-aware cancellation.
package retrier

import (
	"context"
	"time"
)

// clock abstracts time for testing.
type clock interface {
	Now() time.Time
	Sleep(d time.Duration)
}

type realClock struct{}

func (realClock) Now() time.Time          { return time.Now() }
func (realClock) Sleep(d time.Duration)   { time.Sleep(d) }

// Retrier retries a function up to MaxAttempts times with a fixed delay
// between attempts. It respects context cancellation.
type Retrier struct {
	maxAttempts int
	delay       time.Duration
	clk         clock
}

// New returns a Retrier with the given attempt limit and inter-attempt delay.
func New(maxAttempts int, delay time.Duration) *Retrier {
	return newWithClock(maxAttempts, delay, realClock{})
}

func newWithClock(maxAttempts int, delay time.Duration, clk clock) *Retrier {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return &Retrier{maxAttempts: maxAttempts, delay: delay, clk: clk}
}

// Do calls fn up to MaxAttempts times. It stops early if ctx is cancelled
// or fn returns nil. The last non-nil error is returned.
func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	var last error
	for i := 0; i < r.maxAttempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := fn(); err == nil {
			return nil
		} else {
			last = err
		}
		if i < r.maxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-after(ctx, r.delay, r.clk):
			}
		}
	}
	return last
}

// Attempts returns the configured maximum attempt count.
func (r *Retrier) Attempts() int { return r.maxAttempts }

// after returns a channel that fires after d using the provided clock.
// It falls back to a real timer so the select can still honour ctx.Done.
func after(_ context.Context, d time.Duration, clk clock) <-chan struct{} {
	ch := make(chan struct{}, 1)
	go func() {
		clk.Sleep(d)
		ch <- struct{}{}
	}()
	return ch
}
