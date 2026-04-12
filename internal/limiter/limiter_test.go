package limiter

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) clock {
	return func() time.Time { return t }
}

func TestAllow_UnderMax(t *testing.T) {
	l := newWithClock(3, time.Minute, fixedClock(epoch))
	for i := 0; i < 3; i++ {
		if !l.Allow(8080, "tcp") {
			t.Fatalf("call %d: expected Allow=true", i+1)
		}
	}
}

func TestAllow_ExceedsMax(t *testing.T) {
	l := newWithClock(2, time.Minute, fixedClock(epoch))
	l.Allow(9000, "tcp")
	l.Allow(9000, "tcp")
	if l.Allow(9000, "tcp") {
		t.Fatal("expected Allow=false after exceeding max")
	}
}

func TestAllow_WindowResets(t *testing.T) {
	now := epoch
	l := newWithClock(1, time.Minute, func() time.Time { return now })

	if !l.Allow(443, "tcp") {
		t.Fatal("first call should be allowed")
	}
	if l.Allow(443, "tcp") {
		t.Fatal("second call within window should be denied")
	}

	// Advance past the window.
	now = now.Add(2 * time.Minute)
	if !l.Allow(443, "tcp") {
		t.Fatal("first call after window reset should be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := newWithClock(1, time.Minute, fixedClock(epoch))
	if !l.Allow(80, "tcp") {
		t.Fatal("port 80 first call should be allowed")
	}
	if !l.Allow(443, "tcp") {
		t.Fatal("port 443 first call should be allowed (independent key)")
	}
	if l.Allow(80, "tcp") {
		t.Fatal("port 80 second call should be denied")
	}
}

func TestAllow_ProtocolsAreIndependent(t *testing.T) {
	l := newWithClock(1, time.Minute, fixedClock(epoch))
	if !l.Allow(53, "tcp") {
		t.Fatal("tcp:53 should be allowed")
	}
	if !l.Allow(53, "udp") {
		t.Fatal("udp:53 should be allowed (different protocol key)")
	}
}

func TestReset_ClearsState(t *testing.T) {
	l := newWithClock(1, time.Minute, fixedClock(epoch))
	l.Allow(8080, "tcp")
	if l.Allow(8080, "tcp") {
		t.Fatal("should be denied before reset")
	}
	l.Reset()
	if !l.Allow(8080, "tcp") {
		t.Fatal("should be allowed after reset")
	}
}
