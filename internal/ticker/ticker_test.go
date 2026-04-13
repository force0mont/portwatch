package ticker_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ticker"
)

func TestNew_PanicsOnZeroInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero interval")
		}
	}()
	ticker.New(0)
}

func TestRun_CancelsCleanly(t *testing.T) {
	tk := ticker.New(10 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := tk.Run(ctx, func(time.Time) {})
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRun_FiresMultipleTimes(t *testing.T) {
	var count atomic.Int32
	tk := ticker.New(10 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Millisecond)
	defer cancel()

	_ = tk.Run(ctx, func(time.Time) { count.Add(1) })

	got := int(count.Load())
	if got < 3 {
		t.Fatalf("expected at least 3 ticks in 55ms, got %d", got)
	}
}

func TestRun_PassesTimestamp(t *testing.T) {
	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	var received time.Time

	tk := ticker.New(
		5*time.Millisecond,
		ticker.WithJitter(func(d time.Duration) time.Duration { return d }),
	)
	// Patch via unexported option is unavailable from _test package;
	// verify timestamp is non-zero and monotonically recent instead.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	before := time.Now()
	_ = tk.Run(ctx, func(at time.Time) {
		if received.IsZero() {
			received = at
			cancel()
		}
	})

	_ = fixed // suppress unused warning
	if received.Before(before) {
		t.Fatalf("timestamp %v is before test start %v", received, before)
	}
}

func TestRun_WithJitter_UsesJitteredDelay(t *testing.T) {
	var count atomic.Int32
	// jitter that halves the interval so ticks arrive faster
	half := func(d time.Duration) time.Duration { return d / 2 }
	tk := ticker.New(20*time.Millisecond, ticker.WithJitter(half))
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Millisecond)
	defer cancel()

	_ = tk.Run(ctx, func(time.Time) { count.Add(1) })

	got := int(count.Load())
	if got < 4 {
		t.Fatalf("expected >=4 ticks with halved jitter in 55ms, got %d", got)
	}
}
