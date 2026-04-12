package backoff

import (
	"testing"
	"time"
)

type fakeClock struct{ now time.Time }

func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Advance(d time.Duration) { f.now = f.now.Add(d) }

func newFake() (*Backoff, *fakeClock) {
	clk := &fakeClock{now: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	b := newWithClock(100*time.Millisecond, 800*time.Millisecond, clk)
	return b, clk
}

func TestReady_NoFailures_ReturnsTrue(t *testing.T) {
	b, _ := newFake()
	if !b.Ready("svc") {
		t.Fatal("expected ready for unseen key")
	}
}

func TestFailure_BlocksUntilDelay(t *testing.T) {
	b, clk := newFake()
	b.Failure("svc") // delay = 100 ms
	if b.Ready("svc") {
		t.Fatal("should not be ready immediately after failure")
	}
	clk.Advance(100 * time.Millisecond)
	if !b.Ready("svc") {
		t.Fatal("should be ready after delay has elapsed")
	}
}

func TestFailure_ExponentialGrowth(t *testing.T) {
	b, _ := newFake()
	b.Failure("svc") // 1st → 100 ms
	b.Failure("svc") // 2nd → 200 ms
	b.Failure("svc") // 3rd → 400 ms
	if got := b.Failures("svc"); got != 3 {
		t.Fatalf("expected 3 failures, got %d", got)
	}
}

func TestFailure_CapsAtMax(t *testing.T) {
	b, clk := newFake()
	for i := 0; i < 10; i++ {
		b.Failure("svc")
	}
	// Advance by max (800 ms); key must be ready.
	clk.Advance(800 * time.Millisecond)
	if !b.Ready("svc") {
		t.Fatal("expected ready after max delay elapsed")
	}
}

func TestSuccess_ResetsState(t *testing.T) {
	b, _ := newFake()
	b.Failure("svc")
	b.Failure("svc")
	b.Success("svc")
	if !b.Ready("svc") {
		t.Fatal("expected ready after success reset")
	}
	if got := b.Failures("svc"); got != 0 {
		t.Fatalf("expected 0 failures after reset, got %d", got)
	}
}

func TestReady_IndependentKeys(t *testing.T) {
	b, _ := newFake()
	b.Failure("a")
	if !b.Ready("b") {
		t.Fatal("key b should be unaffected by key a failure")
	}
}
