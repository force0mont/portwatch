package aggregator

import (
	"testing"
	"time"

	"github.com/joemiller/portwatch/internal/scanner"
	"github.com/joemiller/portwatch/internal/state"
)

var fixedNow = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeEvent(kind state.EventKind, port uint16) state.PortEvent {
	return state.PortEvent{
		Kind: kind,
		Port: scanner.Port{Port: port, Protocol: "tcp", Address: "0.0.0.0"},
	}
}

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestReady_BeforeWindow_ReturnsFalse(t *testing.T) {
	a := newWithClock(5*time.Second, fixedClock(fixedNow))
	if a.Ready() {
		t.Fatal("expected not ready before window elapses")
	}
}

func TestReady_AfterWindow_ReturnsTrue(t *testing.T) {
	clock := fixedNow
	a := newWithClock(5*time.Second, func() time.Time { return clock })
	clock = fixedNow.Add(6 * time.Second)
	if !a.Ready() {
		t.Fatal("expected ready after window elapses")
	}
}

func TestFlush_EmptyBuffer_ReturnsEmptySummary(t *testing.T) {
	a := New(time.Second)
	s := a.Flush()
	if len(s.Appeared) != 0 || len(s.Disappeared) != 0 {
		t.Fatalf("expected empty summary, got %+v", s)
	}
}

func TestFlush_SegregatesAppearedAndDisappeared(t *testing.T) {
	a := newWithClock(time.Second, fixedClock(fixedNow))
	a.Add(makeEvent(state.EventAppeared, 80))
	a.Add(makeEvent(state.EventAppeared, 443))
	a.Add(makeEvent(state.EventDisappeared, 8080))

	s := a.Flush()
	if len(s.Appeared) != 2 {
		t.Fatalf("expected 2 appeared, got %d", len(s.Appeared))
	}
	if len(s.Disappeared) != 1 {
		t.Fatalf("expected 1 disappeared, got %d", len(s.Disappeared))
	}
}

func TestFlush_ClearsBufferAfterFlush(t *testing.T) {
	a := newWithClock(time.Second, fixedClock(fixedNow))
	a.Add(makeEvent(state.EventAppeared, 22))
	a.Flush()

	s := a.Flush()
	if len(s.Appeared) != 0 {
		t.Fatalf("expected empty buffer after second flush, got %d events", len(s.Appeared))
	}
}

func TestFlush_SetsCollectedAt(t *testing.T) {
	a := newWithClock(time.Second, fixedClock(fixedNow))
	s := a.Flush()
	if !s.CollectedAt.Equal(fixedNow) {
		t.Fatalf("expected CollectedAt=%v, got %v", fixedNow, s.CollectedAt)
	}
}

func TestFlush_ResetsReadyWindow(t *testing.T) {
	clock := fixedNow.Add(10 * time.Second)
	a := newWithClock(5*time.Second, func() time.Time { return clock })
	if !a.Ready() {
		t.Fatal("expected ready before flush")
	}
	a.Flush()
	// After flush the lastFlush is reset to clock; window hasn't elapsed again.
	if a.Ready() {
		t.Fatal("expected not ready immediately after flush")
	}
}
