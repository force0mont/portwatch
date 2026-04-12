package circuitbreaker

import (
	"testing"
	"time"
)

// fixedClock returns a clock whose value can be advanced manually.
func fixedClock(initial time.Time) (*time.Time, clock) {
	t := initial
	return &t, func() time.Time { return t }
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b := New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_OpensAfterMaxFailures(t *testing.T) {
	b := New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_StaysClosedBelowThreshold(t *testing.T) {
	b := New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_RecoveryAfterCooldown(t *testing.T) {
	now := time.Now()
	ptr, clk := fixedClock(now)
	b := newWithClock(2, 5*time.Second, clk)

	b.RecordFailure()
	b.RecordFailure()

	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen before cooldown, got %v", err)
	}

	*ptr = now.Add(6 * time.Second)

	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
}

func TestRecordSuccess_ResetsClosed(t *testing.T) {
	b := New(2, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if !b.IsOpen() {
		t.Fatal("expected circuit to be open")
	}
	// Simulate cooldown then probe success.
	b.mu.Lock()
	b.state = stateClosed
	b.failures = 0
	b.mu.Unlock()
	b.RecordSuccess()
	if b.IsOpen() {
		t.Fatal("expected circuit to be closed after success")
	}
}

func TestIsOpen_ReflectsState(t *testing.T) {
	b := New(1, time.Second)
	if b.IsOpen() {
		t.Fatal("should not be open initially")
	}
	b.RecordFailure()
	if !b.IsOpen() {
		t.Fatal("should be open after threshold failure")
	}
}

func TestAllow_RemainsOpenWithinCooldown(t *testing.T) {
	now := time.Now()
	ptr, clk := fixedClock(now)
	b := newWithClock(1, 10*time.Second, clk)

	b.RecordFailure()
	*ptr = now.Add(5 * time.Second)

	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen within cooldown, got %v", err)
	}
}
