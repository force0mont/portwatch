package debounce

import (
	"testing"
	"time"

	"github.com/rjbrown57/portwatch/internal/state"
)

// fakeClock lets tests control the current time.
type fakeClock struct{ now time.Time }

func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Advance(d time.Duration) { f.now = f.now.Add(d) }

func makeEvent(proto, addr string) state.Event {
	return state.Event{Protocol: proto, Addr: addr, Kind: state.Appeared}
}

func TestFeed_EventBelowWindow_NotReleased(t *testing.T) {
	clk := &fakeClock{now: time.Now()}
	d := newWithClock(2*time.Second, clk)

	events := []state.Event{makeEvent("tcp", "0.0.0.0:8080")}

	// First feed — starts the clock.
	out := d.Feed(events)
	if len(out) != 0 {
		t.Fatalf("expected 0 stable events, got %d", len(out))
	}
	if d.Len() != 1 {
		t.Fatalf("expected 1 pending event, got %d", d.Len())
	}

	// Advance only 1 s — still below the 2 s window.
	clk.Advance(1 * time.Second)
	out = d.Feed(events)
	if len(out) != 0 {
		t.Fatalf("expected 0 stable events after 1 s, got %d", len(out))
	}
}

func TestFeed_EventAboveWindow_Released(t *testing.T) {
	clk := &fakeClock{now: time.Now()}
	d := newWithClock(2*time.Second, clk)

	events := []state.Event{makeEvent("tcp", "0.0.0.0:9090")}

	d.Feed(events)
	clk.Advance(3 * time.Second)

	out := d.Feed(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 stable event, got %d", len(out))
	}
	if d.Len() != 0 {
		t.Fatalf("expected pending set to be empty after release, got %d", d.Len())
	}
}

func TestFeed_TransientEvent_Evicted(t *testing.T) {
	clk := &fakeClock{now: time.Now()}
	d := newWithClock(2*time.Second, clk)

	events := []state.Event{makeEvent("udp", "0.0.0.0:5353")}
	d.Feed(events)

	// Next feed omits the event — it should be evicted.
	out := d.Feed(nil)
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
	if d.Len() != 0 {
		t.Fatalf("expected empty pending set, got %d", d.Len())
	}
}

func TestFeed_MultipleEvents_IndependentWindows(t *testing.T) {
	clk := &fakeClock{now: time.Now()}
	d := newWithClock(2*time.Second, clk)

	e1 := makeEvent("tcp", "0.0.0.0:80")
	e2 := makeEvent("tcp", "0.0.0.0:443")

	d.Feed([]state.Event{e1, e2})
	clk.Advance(3 * time.Second)

	// e2 disappears before the window expires for it (already advanced 3 s, so
	// both would be stable — feed only e1 to simulate e2 going away first).
	out := d.Feed([]state.Event{e1})
	if len(out) != 1 {
		t.Fatalf("expected 1 stable event, got %d", len(out))
	}
	if out[0].Addr != e1.Addr {
		t.Errorf("expected %s, got %s", e1.Addr, out[0].Addr)
	}
}

func TestNew_DefaultWindow(t *testing.T) {
	d := New(5 * time.Second)
	if d.Window != 5*time.Second {
		t.Errorf("unexpected window: %v", d.Window)
	}
}
