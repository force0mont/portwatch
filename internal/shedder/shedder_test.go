package shedder

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_HighPriority_NeverShed(t *testing.T) {
	s := newWithClock(time.Second, 0, 0, fixedClock(epoch))
	for i := 0; i < 100; i++ {
		if !s.Allow(PriorityHigh) {
			t.Fatal("high priority should never be shed")
		}
	}
}

func TestAllow_LowPriority_UnderMax(t *testing.T) {
	s := newWithClock(time.Second, 3, 10, fixedClock(epoch))
	for i := 0; i < 3; i++ {
		if !s.Allow(PriorityLow) {
			t.Fatalf("call %d should be allowed", i)
		}
	}
}

func TestAllow_LowPriority_ExceedsMax(t *testing.T) {
	s := newWithClock(time.Second, 2, 10, fixedClock(epoch))
	s.Allow(PriorityLow)
	s.Allow(PriorityLow)
	if s.Allow(PriorityLow) {
		t.Fatal("third low-priority call should be shed")
	}
}

func TestAllow_NormalPriority_ExceedsMax(t *testing.T) {
	s := newWithClock(time.Second, 10, 2, fixedClock(epoch))
	s.Allow(PriorityNormal)
	s.Allow(PriorityNormal)
	if s.Allow(PriorityNormal) {
		t.Fatal("third normal-priority call should be shed")
	}
}

func TestAllow_WindowResets_AllowsAgain(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	s := newWithClock(time.Second, 1, 10, clock)
	s.Allow(PriorityLow) // fills the window
	if s.Allow(PriorityLow) {
		t.Fatal("should be shed within window")
	}
	now = now.Add(2 * time.Second) // advance past window
	if !s.Allow(PriorityLow) {
		t.Fatal("should be allowed after window reset")
	}
}

func TestAllow_IndependentPriorities(t *testing.T) {
	s := newWithClock(time.Second, 1, 1, fixedClock(epoch))
	s.Allow(PriorityLow)
	// low is exhausted, but normal should still be independent
	if !s.Allow(PriorityNormal) {
		t.Fatal("normal priority should be independent of low")
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	s := newWithClock(time.Second, 1, 1, fixedClock(epoch))
	s.Allow(PriorityLow)
	s.Allow(PriorityNormal)
	s.Reset()
	if !s.Allow(PriorityLow) {
		t.Fatal("low should be allowed after reset")
	}
	if !s.Allow(PriorityNormal) {
		t.Fatal("normal should be allowed after reset")
	}
}
