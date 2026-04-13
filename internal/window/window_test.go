package window

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestAdd_ReturnsRunningCount(t *testing.T) {
	w := newWithClock(time.Minute, fixedClock(epoch))
	if got := w.Add("k"); got != 1 {
		t.Fatalf("want 1, got %d", got)
	}
	if got := w.Add("k"); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
}

func TestAdd_IndependentKeys(t *testing.T) {
	w := newWithClock(time.Minute, fixedClock(epoch))
	w.Add("a")
	w.Add("a")
	w.Add("b")
	if got := w.Count("a"); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
	if got := w.Count("b"); got != 1 {
		t.Fatalf("want 1, got %d", got)
	}
}

func TestAdd_EvictsExpiredBuckets(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	w := newWithClock(time.Minute, clock)

	w.Add("k") // recorded at epoch

	now = epoch.Add(90 * time.Second) // advance past window
	if got := w.Add("k"); got != 1 {
		t.Fatalf("old bucket should have been evicted; want 1, got %d", got)
	}
}

func TestCount_DoesNotMutate(t *testing.T) {
	w := newWithClock(time.Minute, fixedClock(epoch))
	w.Add("k")
	w.Count("k")
	w.Count("k")
	if got := w.Count("k"); got != 1 {
		t.Fatalf("Count should not add entries; want 1, got %d", got)
	}
}

func TestCount_UnknownKey_ReturnsZero(t *testing.T) {
	w := New(time.Minute)
	if got := w.Count("missing"); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}
}

func TestReset_ClearsKey(t *testing.T) {
	w := newWithClock(time.Minute, fixedClock(epoch))
	w.Add("k")
	w.Add("k")
	w.Reset("k")
	if got := w.Count("k"); got != 0 {
		t.Fatalf("want 0 after reset, got %d", got)
	}
}

func TestReset_OtherKeyUnaffected(t *testing.T) {
	w := newWithClock(time.Minute, fixedClock(epoch))
	w.Add("a")
	w.Add("b")
	w.Reset("a")
	if got := w.Count("b"); got != 1 {
		t.Fatalf("want 1, got %d", got)
	}
}
