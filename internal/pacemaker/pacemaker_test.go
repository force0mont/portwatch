package pacemaker

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestBeat_FirstCall_ReturnsTrue(t *testing.T) {
	p := newWithClock(time.Second, fixedClock(epoch))
	if !p.Beat() {
		t.Fatal("expected first beat to return true")
	}
}

func TestBeat_WithinThreshold_ReturnsTrue(t *testing.T) {
	now := epoch
	p := newWithClock(time.Second, func() time.Time { return now })
	p.Beat()
	now = epoch.Add(500 * time.Millisecond)
	if !p.Beat() {
		t.Fatal("expected beat within threshold to return true")
	}
	if p.Missed() != 0 {
		t.Fatalf("expected 0 missed, got %d", p.Missed())
	}
}

func TestBeat_ExceedsThreshold_ReturnsFalse(t *testing.T) {
	now := epoch
	p := newWithClock(time.Second, func() time.Time { return now })
	p.Beat()
	now = epoch.Add(2 * time.Second)
	if p.Beat() {
		t.Fatal("expected beat beyond threshold to return false")
	}
	if p.Missed() != 1 {
		t.Fatalf("expected 1 missed, got %d", p.Missed())
	}
}

func TestMissed_AccumulatesAcrossBeats(t *testing.T) {
	now := epoch
	p := newWithClock(time.Second, func() time.Time { return now })
	p.Beat()
	for i := 1; i <= 3; i++ {
		now = epoch.Add(time.Duration(i) * 2 * time.Second)
		p.Beat()
	}
	if p.Missed() != 3 {
		t.Fatalf("expected 3 missed, got %d", p.Missed())
	}
}

func TestReset_ClearsMissedAndLastBeat(t *testing.T) {
	now := epoch
	p := newWithClock(time.Second, func() time.Time { return now })
	p.Beat()
	now = epoch.Add(5 * time.Second)
	p.Beat()
	p.Reset()
	if p.Missed() != 0 {
		t.Fatalf("expected 0 missed after reset, got %d", p.Missed())
	}
	if !p.LastBeat().IsZero() {
		t.Fatal("expected zero LastBeat after reset")
	}
}

func TestLastBeat_ReturnsCurrentTime(t *testing.T) {
	p := newWithClock(time.Second, fixedClock(epoch))
	p.Beat()
	if !p.LastBeat().Equal(epoch) {
		t.Fatalf("expected %v, got %v", epoch, p.LastBeat())
	}
}
