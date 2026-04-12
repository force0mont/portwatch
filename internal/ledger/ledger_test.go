package ledger

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedNow }

func TestRecordAppeared_CreatesEntry(t *testing.T) {
	l := newWithClock(fixedClock)
	l.RecordAppeared(8080, "tcp")

	e, ok := l.Get(8080, "tcp")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Appeared != 1 {
		t.Fatalf("expected Appeared=1, got %d", e.Appeared)
	}
	if e.Disappeared != 0 {
		t.Fatalf("expected Disappeared=0, got %d", e.Disappeared)
	}
	if !e.LastSeen.Equal(fixedNow) {
		t.Fatalf("unexpected LastSeen: %v", e.LastSeen)
	}
}

func TestRecordDisappeared_IncrementsCounter(t *testing.T) {
	l := newWithClock(fixedClock)
	l.RecordDisappeared(443, "tcp")
	l.RecordDisappeared(443, "tcp")

	e, ok := l.Get(443, "tcp")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Disappeared != 2 {
		t.Fatalf("expected Disappeared=2, got %d", e.Disappeared)
	}
}

func TestGet_UnknownPort_ReturnsFalse(t *testing.T) {
	l := New()
	_, ok := l.Get(9999, "tcp")
	if ok {
		t.Fatal("expected ok=false for unknown port")
	}
}

func TestAll_ReturnsCopyOfEntries(t *testing.T) {
	l := newWithClock(fixedClock)
	l.RecordAppeared(80, "tcp")
	l.RecordAppeared(53, "udp")

	all := l.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestReset_ClearsAllEntries(t *testing.T) {
	l := newWithClock(fixedClock)
	l.RecordAppeared(80, "tcp")
	l.Reset()

	all := l.All()
	if len(all) != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", len(all))
	}
}

func TestProtocol_Independent_Keys(t *testing.T) {
	l := newWithClock(fixedClock)
	l.RecordAppeared(53, "tcp")
	l.RecordAppeared(53, "udp")

	tcp, _ := l.Get(53, "tcp")
	udp, _ := l.Get(53, "udp")

	if tcp.Key == udp.Key {
		t.Fatal("tcp and udp should have distinct keys")
	}
	if tcp.Appeared != 1 || udp.Appeared != 1 {
		t.Fatal("each protocol should have Appeared=1")
	}
}
