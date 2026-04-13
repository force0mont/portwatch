package budget

import (
	"testing"
	"time"
)

var (
	baseTime  = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	fixedNow  = baseTime
)

func fixedClock() time.Time { return fixedNow }

func advanceBy(d time.Duration) { fixedNow = fixedNow.Add(d) }

func resetClock() { fixedNow = baseTime }

func TestAllow_UnderMax(t *testing.T) {
	resetClock()
	b := newWithClock(3, time.Minute, fixedClock)
	for i := 0; i < 3; i++ {
		if !b.Allow("k") {
			t.Fatalf("expected Allow true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsMax(t *testing.T) {
	resetClock()
	b := newWithClock(2, time.Minute, fixedClock)
	b.Allow("k")
	b.Allow("k")
	if b.Allow("k") {
		t.Fatal("expected Allow false after budget exhausted")
	}
}

func TestAllow_WindowResets(t *testing.T) {
	resetClock()
	b := newWithClock(1, time.Minute, fixedClock)
	b.Allow("k")
	if b.Allow("k") {
		t.Fatal("expected false within window")
	}
	advanceBy(61 * time.Second)
	if !b.Allow("k") {
		t.Fatal("expected true after window reset")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	resetClock()
	b := newWithClock(1, time.Minute, fixedClock)
	b.Allow("a")
	if !b.Allow("b") {
		t.Fatal("key b should be independent of key a")
	}
}

func TestRemaining_DecreasesWithUse(t *testing.T) {
	resetClock()
	b := newWithClock(3, time.Minute, fixedClock)
	if got := b.Remaining("x"); got != 3 {
		t.Fatalf("want 3 remaining, got %d", got)
	}
	b.Allow("x")
	if got := b.Remaining("x"); got != 2 {
		t.Fatalf("want 2 remaining, got %d", got)
	}
}

func TestReset_ClearsState(t *testing.T) {
	resetClock()
	b := newWithClock(1, time.Minute, fixedClock)
	b.Allow("z")
	if b.Allow("z") {
		t.Fatal("expected false before reset")
	}
	b.Reset("z")
	if !b.Allow("z") {
		t.Fatal("expected true after reset")
	}
}

func TestNew_DefaultsApplied(t *testing.T) {
	b := newWithClock(0, 0, fixedClock)
	if b.max != 1 {
		t.Fatalf("want max=1 for zero input, got %d", b.max)
	}
	if b.window != time.Minute {
		t.Fatalf("want window=1m for zero input, got %v", b.window)
	}
}
