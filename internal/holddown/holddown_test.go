package holddown

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestSeen_FirstCall_ReturnsTrue(t *testing.T) {
	h := newWithClock(5*time.Second, fixedClock(epoch))
	if !h.Seen(8080, "tcp") {
		t.Fatal("expected true on first Seen")
	}
}

func TestSeen_SecondCall_ReturnsFalse(t *testing.T) {
	h := newWithClock(5*time.Second, fixedClock(epoch))
	h.Seen(8080, "tcp")
	if h.Seen(8080, "tcp") {
		t.Fatal("expected false on second Seen")
	}
}

func TestSeen_AfterGone_ReturnsTrue(t *testing.T) {
	h := newWithClock(5*time.Second, fixedClock(epoch))
	h.Seen(8080, "tcp")
	h.Gone(8080, "tcp")
	if !h.Seen(8080, "tcp") {
		t.Fatal("expected true after Gone")
	}
}

func TestSeen_DifferentProtocols_Independent(t *testing.T) {
	h := newWithClock(5*time.Second, fixedClock(epoch))
	if !h.Seen(53, "tcp") {
		t.Fatal("tcp:53 first seen should return true")
	}
	if !h.Seen(53, "udp") {
		t.Fatal("udp:53 first seen should return true")
	}
}

func TestPrune_RemovesStaleEntries(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	h := newWithClock(5*time.Second, clock)

	h.Seen(9000, "tcp")
	if h.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", h.Len())
	}

	// advance past quiet period
	now = epoch.Add(10 * time.Second)
	removed := h.Prune()
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if h.Len() != 0 {
		t.Fatalf("expected 0 entries after prune, got %d", h.Len())
	}
}

func TestPrune_KeepsRecentEntries(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	h := newWithClock(30*time.Second, clock)

	h.Seen(443, "tcp")
	now = epoch.Add(5 * time.Second)
	removed := h.Prune()
	if removed != 0 {
		t.Fatalf("expected 0 removed, got %d", removed)
	}
	if h.Len() != 1 {
		t.Fatalf("expected 1 entry retained, got %d", h.Len())
	}
}

func TestPrune_AfterPrune_SeenFiresAgain(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	h := newWithClock(5*time.Second, clock)

	h.Seen(8443, "tcp")
	now = epoch.Add(10 * time.Second)
	h.Prune()

	if !h.Seen(8443, "tcp") {
		t.Fatal("expected true after prune cleared the entry")
	}
}
