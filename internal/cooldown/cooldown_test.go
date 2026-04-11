package cooldown

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) clock {
	return func() time.Time { return t }
}

func TestReady_UnseenKey_ReturnsTrue(t *testing.T) {
	tr := newWithClock(5*time.Second, fixedClock(epoch))
	if !tr.Ready("port:tcp:8080") {
		t.Fatal("expected Ready to return true for unseen key")
	}
}

func TestReady_WithinCooldown_ReturnsFalse(t *testing.T) {
	current := epoch
	tr := newWithClock(5*time.Second, func() time.Time { return current })

	tr.Mark("k")
	current = epoch.Add(2 * time.Second)

	if tr.Ready("k") {
		t.Fatal("expected Ready to return false within cooldown")
	}
}

func TestReady_AfterCooldown_ReturnsTrue(t *testing.T) {
	current := epoch
	tr := newWithClock(5*time.Second, func() time.Time { return current })

	tr.Mark("k")
	current = epoch.Add(6 * time.Second)

	if !tr.Ready("k") {
		t.Fatal("expected Ready to return true after cooldown elapsed")
	}
}

func TestReadyAndMark_FirstCall_ReturnsTrue(t *testing.T) {
	tr := newWithClock(5*time.Second, fixedClock(epoch))
	if !tr.ReadyAndMark("k") {
		t.Fatal("expected ReadyAndMark to return true on first call")
	}
}

func TestReadyAndMark_SecondCallWithinCooldown_ReturnsFalse(t *testing.T) {
	current := epoch
	tr := newWithClock(5*time.Second, func() time.Time { return current })

	tr.ReadyAndMark("k")
	current = epoch.Add(1 * time.Second)

	if tr.ReadyAndMark("k") {
		t.Fatal("expected ReadyAndMark to return false within cooldown")
	}
}

func TestReset_AllowsImmediateReady(t *testing.T) {
	current := epoch
	tr := newWithClock(5*time.Second, func() time.Time { return current })

	tr.Mark("k")
	tr.Reset("k")

	if !tr.Ready("k") {
		t.Fatal("expected Ready to return true after Reset")
	}
}

func TestReady_IndependentKeys(t *testing.T) {
	current := epoch
	tr := newWithClock(5*time.Second, func() time.Time { return current })

	tr.Mark("a")
	current = epoch.Add(2 * time.Second)

	if tr.Ready("a") {
		t.Fatal("key 'a' should still be in cooldown")
	}
	if !tr.Ready("b") {
		t.Fatal("key 'b' should be ready (never marked)")
	}
}
