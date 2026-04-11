package suppress

import (
	"testing"
	"time"
)

var (
	testPort  uint16 = 8080
	testProto        = "tcp"
)

func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestIsSuppressed_FirstCall_NotSuppressed(t *testing.T) {
	now := time.Now()
	s := newWithClock(30*time.Second, fixedClock(now))

	if s.IsSuppressed(testPort, testProto) {
		t.Fatal("expected first call to not be suppressed")
	}
}

func TestIsSuppressed_SecondCall_WithinCooldown_Suppressed(t *testing.T) {
	now := time.Now()
	s := newWithClock(30*time.Second, fixedClock(now))

	s.IsSuppressed(testPort, testProto) // record

	if !s.IsSuppressed(testPort, testProto) {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestIsSuppressed_AfterCooldown_NotSuppressed(t *testing.T) {
	now := time.Now()
	s := newWithClock(30*time.Second, fixedClock(now))

	s.IsSuppressed(testPort, testProto) // record at 'now'

	// advance clock beyond cooldown
	s.clock = fixedClock(now.Add(31 * time.Second))

	if s.IsSuppressed(testPort, testProto) {
		t.Fatal("expected call after cooldown to not be suppressed")
	}
}

func TestIsSuppressed_DifferentProtocols_Independent(t *testing.T) {
	now := time.Now()
	s := newWithClock(30*time.Second, fixedClock(now))

	s.IsSuppressed(testPort, "tcp")

	if s.IsSuppressed(testPort, "udp") {
		t.Fatal("udp should not be suppressed when only tcp was recorded")
	}
}

func TestReset_AllowsImmediateAlert(t *testing.T) {
	now := time.Now()
	s := newWithClock(30*time.Second, fixedClock(now))

	s.IsSuppressed(testPort, testProto) // record
	s.Reset(testPort, testProto)

	if s.IsSuppressed(testPort, testProto) {
		t.Fatal("expected alert after Reset to not be suppressed")
	}
}

func TestLen_TracksEntries(t *testing.T) {
	now := time.Now()
	s := newWithClock(30*time.Second, fixedClock(now))

	if s.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", s.Len())
	}

	s.IsSuppressed(8080, "tcp")
	s.IsSuppressed(9090, "tcp")

	if s.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", s.Len())
	}

	s.Reset(8080, "tcp")

	if s.Len() != 1 {
		t.Fatalf("expected 1 entry after reset, got %d", s.Len())
	}
}
