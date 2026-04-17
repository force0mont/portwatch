package fence

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_BelowThreshold_ReturnsFalse(t *testing.T) {
	f := newWithClock(time.Minute, 3, fixedClock(epoch))
	_, tripped := f.Record("k")
	if tripped {
		t.Fatal("expected false below threshold")
	}
}

func TestRecord_AtThreshold_ReturnsTrue(t *testing.T) {
	f := newWithClock(time.Minute, 3, fixedClock(epoch))
	f.Record("k")
	f.Record("k")
	trip, ok := f.Record("k")
	if !ok {
		t.Fatal("expected tripped at threshold")
	}
	if trip.Key != "k" {
		t.Fatalf("unexpected key %q", trip.Key)
	}
	if trip.Count != 3 {
		t.Fatalf("expected count 3, got %d", trip.Count)
	}
}

func TestRecord_IndependentKeys(t *testing.T) {
	f := newWithClock(time.Minute, 2, fixedClock(epoch))
	f.Record("a")
	_, ok := f.Record("b")
	if ok {
		t.Fatal("key b should not trip from key a's events")
	}
}

func TestRecord_WindowEvictsOldEvents(t *testing.T) {
	now := epoch
	f := newWithClock(time.Second*10, 2, func() time.Time { return now })
	f.Record("k") // t=0
	now = epoch.Add(time.Second * 11)
	// first event is now outside the window; single new event should not trip
	_, ok := f.Record("k")
	if ok {
		t.Fatal("stale event should have been evicted")
	}
}

func TestReset_ClearsHistory(t *testing.T) {
	f := newWithClock(time.Minute, 2, fixedClock(epoch))
	f.Record("k")
	f.Reset("k")
	_, ok := f.Record("k")
	if ok {
		t.Fatal("expected false after reset")
	}
}
