package throttle

import (
	"testing"
	"time"
)

var (
	now   = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	fixed = func() time.Time { return now }
)

func TestAllow_UnderMax(t *testing.T) {
	th := newWithClock(3, time.Minute, fixed)
	for i := 0; i < 3; i++ {
		if !th.Allow(8080, "tcp") {
			t.Fatalf("call %d: expected Allow=true", i+1)
		}
	}
}

func TestAllow_ExceedsMax(t *testing.T) {
	th := newWithClock(2, time.Minute, fixed)
	th.Allow(443, "tcp")
	th.Allow(443, "tcp")
	if th.Allow(443, "tcp") {
		t.Fatal("expected Allow=false after exceeding max")
	}
}

func TestAllow_WindowResets(t *testing.T) {
	current := now
	th := newWithClock(1, time.Minute, func() time.Time { return current })

	th.Allow(22, "tcp") // count=1, allowed
	th.Allow(22, "tcp") // count=2, throttled

	current = now.Add(2 * time.Minute) // advance past window
	if !th.Allow(22, "tcp") {
		t.Fatal("expected Allow=true after window reset")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th := newWithClock(1, time.Minute, fixed)
	th.Allow(80, "tcp")  // fills tcp:80
	if !th.Allow(80, "udp") {
		t.Fatal("udp:80 should be independent of tcp:80")
	}
}

func TestReset_ClearsState(t *testing.T) {
	th := newWithClock(1, time.Minute, fixed)
	th.Allow(9000, "tcp")
	th.Allow(9000, "tcp") // now throttled
	th.Reset()
	if !th.Allow(9000, "tcp") {
		t.Fatal("expected Allow=true after Reset")
	}
}

func TestNew_NotNil(t *testing.T) {
	th := New(5, 30*time.Second)
	if th == nil {
		t.Fatal("expected non-nil Throttle")
	}
}
