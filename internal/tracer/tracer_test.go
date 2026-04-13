package tracer

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

var (
	t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 = t0.Add(5 * time.Second)
	t2 = t0.Add(15 * time.Second)
)

func fixedClock(ts time.Time) func() time.Time {
	return func() time.Time { return ts }
}

func makePort(proto string, port int) scanner.Port {
	return scanner.Port{Protocol: proto, Address: "0.0.0.0", Port: port}
}

func TestObserve_NewPort_SetsFirstSeen(t *testing.T) {
	tr := newWithClock(fixedClock(t0))
	e := tr.Observe(makePort("tcp", 8080))
	if !e.FirstSeen.Equal(t0) {
		t.Fatalf("expected FirstSeen=%v got %v", t0, e.FirstSeen)
	}
	if e.Duration != 0 {
		t.Fatalf("expected zero duration on first observe, got %v", e.Duration)
	}
}

func TestObserve_SecondCall_UpdatesDuration(t *testing.T) {
	clock := t0
	tr := newWithClock(func() time.Time { return clock })
	p := makePort("tcp", 9000)
	tr.Observe(p)
	clock = t1
	e := tr.Observe(p)
	if e.Duration != 5*time.Second {
		t.Fatalf("expected 5s duration, got %v", e.Duration)
	}
	if !e.FirstSeen.Equal(t0) {
		t.Fatalf("FirstSeen should not change: got %v", e.FirstSeen)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	tr := newWithClock(fixedClock(t0))
	p := makePort("udp", 53)
	tr.Observe(p)
	tr.Remove(p)
	if _, ok := tr.Get(p); ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestGet_UnknownPort_ReturnsFalse(t *testing.T) {
	tr := New()
	if _, ok := tr.Get(makePort("tcp", 1234)); ok {
		t.Fatal("expected false for unknown port")
	}
}

func TestLen_ReflectsActiveEntries(t *testing.T) {
	tr := newWithClock(fixedClock(t0))
	if tr.Len() != 0 {
		t.Fatal("expected 0")
	}
	tr.Observe(makePort("tcp", 80))
	tr.Observe(makePort("tcp", 443))
	if tr.Len() != 2 {
		t.Fatalf("expected 2, got %d", tr.Len())
	}
	tr.Remove(makePort("tcp", 80))
	if tr.Len() != 1 {
		t.Fatalf("expected 1 after remove, got %d", tr.Len())
	}
}

func TestObserve_DifferentProtocols_Independent(t *testing.T) {
	tr := newWithClock(fixedClock(t0))
	tr.Observe(scanner.Port{Protocol: "tcp", Address: "0.0.0.0", Port: 53})
	tr.Observe(scanner.Port{Protocol: "udp", Address: "0.0.0.0", Port: 53})
	if tr.Len() != 2 {
		t.Fatalf("expected 2 independent entries, got %d", tr.Len())
	}
}
