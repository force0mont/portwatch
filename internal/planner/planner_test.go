package planner

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestNext_ReturnsBaseIntervalWhenNoJitter(t *testing.T) {
	p := newWithClock(10*time.Second, 0, fixedClock(epoch))
	got := p.Next()
	if got != 10*time.Second {
		t.Fatalf("expected 10s, got %v", got)
	}
}

func TestNext_WithJitter_WithinBounds(t *testing.T) {
	interval := 10 * time.Second
	jitter := 2 * time.Second
	p := newWithClock(interval, jitter, fixedClock(epoch))
	got := p.Next()
	if got < interval-jitter/2 || got > interval+jitter/2 {
		t.Fatalf("jittered value %v out of expected range", got)
	}
}

func TestNext_JitterClampedWhenTooLarge(t *testing.T) {
	// jitter >= interval should be clamped to zero
	p := New(5*time.Second, 10*time.Second)
	got := p.Next()
	if got != 5*time.Second {
		t.Fatalf("expected interval unchanged, got %v", got)
	}
}

func TestMark_NoMissedOnFirstCall(t *testing.T) {
	p := newWithClock(10*time.Second, 0, fixedClock(epoch))
	p.Mark(epoch)
	if p.Missed() != 0 {
		t.Fatalf("expected 0 missed, got %d", p.Missed())
	}
}

func TestMark_MissedWhenGapExceedsThreshold(t *testing.T) {
	p := newWithClock(10*time.Second, 0, fixedClock(epoch))
	p.Mark(epoch)
	// Gap of 20s > 1.5 × 10s = 15s → should count as missed
	p.Mark(epoch.Add(20 * time.Second))
	if p.Missed() != 1 {
		t.Fatalf("expected 1 missed, got %d", p.Missed())
	}
}

func TestMark_NotMissedWhenGapWithinThreshold(t *testing.T) {
	p := newWithClock(10*time.Second, 0, fixedClock(epoch))
	p.Mark(epoch)
	p.Mark(epoch.Add(12 * time.Second))
	if p.Missed() != 0 {
		t.Fatalf("expected 0 missed, got %d", p.Missed())
	}
}

func TestReset_ClearsMissedAndLastFire(t *testing.T) {
	p := newWithClock(10*time.Second, 0, fixedClock(epoch))
	p.Mark(epoch)
	p.Mark(epoch.Add(30 * time.Second))
	p.Reset()
	if p.Missed() != 0 {
		t.Fatalf("expected 0 after reset, got %d", p.Missed())
	}
	// After reset the next Mark should not count as missed
	p.Mark(epoch.Add(60 * time.Second))
	if p.Missed() != 0 {
		t.Fatalf("expected 0 after first mark post-reset, got %d", p.Missed())
	}
}
