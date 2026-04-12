package retrier

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock records sleep calls without blocking.
type fakeClock struct {
	slept []time.Duration
}

func (f *fakeClock) Now() time.Time              { return time.Time{} }
func (f *fakeClock) Sleep(d time.Duration)        { f.slept = append(f.slept, d) }

var errBoom = errors.New("boom")

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	r := newWithClock(3, 0, &fakeClock{})
	var calls int
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesUpToMax(t *testing.T) {
	clk := &fakeClock{}
	r := newWithClock(3, time.Millisecond, clk)
	var calls int
	err := r.Do(context.Background(), func() error {
		calls++
		return errBoom
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected errBoom, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	if len(clk.slept) != 2 {
		t.Fatalf("expected 2 sleeps, got %d", len(clk.slept))
	}
}

func TestDo_SucceedsOnRetry(t *testing.T) {
	r := newWithClock(3, 0, &fakeClock{})
	var calls int
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return errBoom
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancelled_StopsEarly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := newWithClock(5, 0, &fakeClock{})
	var calls int32
	err := r.Do(ctx, func() error {
		atomic.AddInt32(&calls, 1)
		return errBoom
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected Canceled, got %v", err)
	}
	if atomic.LoadInt32(&calls) > 0 {
		t.Fatalf("expected 0 calls after pre-cancelled ctx, got %d", calls)
	}
}

func TestAttempts_ReturnsConfiguredValue(t *testing.T) {
	r := New(7, time.Second)
	if r.Attempts() != 7 {
		t.Fatalf("expected 7, got %d", r.Attempts())
	}
}

func TestNew_ClampsZeroAttempts(t *testing.T) {
	r := New(0, 0)
	if r.Attempts() != 1 {
		t.Fatalf("expected 1, got %d", r.Attempts())
	}
}
