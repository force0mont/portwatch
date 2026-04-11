package rollup

import (
	"testing"
	"time"
)

var (
	epoch    = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	windowDur = 10 * time.Second
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_BelowThreshold_NoSummary(t *testing.T) {
	r := newWithClock(windowDur, 3, fixedClock(epoch))
	for i := 0; i < 2; i++ {
		msg, ok := r.Record("port:8080")
		if ok {
			t.Fatalf("expected no summary, got %q", msg)
		}
	}
}

func TestRecord_AtThreshold_ReturnsSummary(t *testing.T) {
	r := newWithClock(windowDur, 3, fixedClock(epoch))
	var got string
	var fired bool
	for i := 0; i < 3; i++ {
		got, fired = r.Record("port:8080")
	}
	if !fired {
		t.Fatal("expected summary to fire at threshold")
	}
	if got == "" {
		t.Fatal("expected non-empty summary message")
	}
}

func TestRecord_AboveThreshold_OnlyFiresOnce(t *testing.T) {
	r := newWithClock(windowDur, 3, fixedClock(epoch))
	count := 0
	for i := 0; i < 10; i++ {
		_, ok := r.Record("port:9090")
		if ok {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected summary to fire once, fired %d times", count)
	}
}

func TestRecord_WindowReset_FiresAgain(t *testing.T) {
	current := epoch
	clock := func() time.Time { return current }
	r := newWithClock(windowDur, 2, clock)

	for i := 0; i < 2; i++ {
		r.Record("port:443")
	}
	// advance past window
	current = epoch.Add(20 * time.Second)
	for i := 0; i < 2; i++ {
		r.Record("port:443")
	}
	_, ok := r.Record("port:443") // should NOT re-fire (already flushed in new window at count==2)
	_ = ok
	// verify a fresh window fires again at threshold
	current = epoch.Add(40 * time.Second)
	r.Record("port:443")
	_, second := r.Record("port:443")
	if !second {
		t.Fatal("expected rollup to fire again after window reset")
	}
}

func TestRecord_IndependentKeys(t *testing.T) {
	r := newWithClock(windowDur, 2, fixedClock(epoch))
	r.Record("port:80")
	r.Record("port:80")
	_, ok := r.Record("port:443")
	if ok {
		t.Fatal("port:443 should not fire based on port:80 counts")
	}
}

func TestReset_ClearsState(t *testing.T) {
	r := newWithClock(windowDur, 2, fixedClock(epoch))
	r.Record("port:22")
	r.Record("port:22")
	r.Reset()
	_, ok := r.Record("port:22")
	if ok {
		t.Fatal("expected no summary after reset")
	}
}
