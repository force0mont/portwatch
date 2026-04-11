package history_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
	"github.com/yourorg/portwatch/internal/state"
)

func makeEvents(kinds ...state.EventKind) []state.Event {
	events := make([]state.Event, len(kinds))
	for i, k := range kinds {
		events[i] = state.Event{Kind: k, Port: uint16(8000 + i), Proto: "tcp"}
	}
	return events
}

func TestNew_DefaultCapacity(t *testing.T) {
	h := history.New(0)
	if h == nil {
		t.Fatal("expected non-nil History")
	}
}

func TestRecord_And_Len(t *testing.T) {
	h := history.New(10)
	h.Record(makeEvents(state.EventAppeared, state.EventAppeared))
	if got := h.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
}

func TestRecord_EvictsOldestWhenFull(t *testing.T) {
	h := history.New(3)
	h.Record(makeEvents(state.EventAppeared, state.EventAppeared, state.EventAppeared))
	// Adding one more should evict the first.
	h.Record(makeEvents(state.EventDisappeared))
	if got := h.Len(); got != 3 {
		t.Fatalf("Len() = %d, want 3 (capacity capped)", got)
	}
	all := h.All()
	if all[2].Event.Kind != state.EventDisappeared {
		t.Errorf("last entry kind = %v, want EventDisappeared", all[2].Event.Kind)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	h := history.New(5)
	h.Record(makeEvents(state.EventAppeared))
	a := h.All()
	a[0].Event.Port = 9999 // mutate copy
	if h.All()[0].Event.Port == 9999 {
		t.Error("All() returned a reference, not a copy")
	}
}

func TestSince_FiltersCorrectly(t *testing.T) {
	h := history.New(10)
	h.Record(makeEvents(state.EventAppeared))
	cutoff := time.Now().UTC()
	time.Sleep(2 * time.Millisecond)
	h.Record(makeEvents(state.EventDisappeared))

	recent := h.Since(cutoff)
	if len(recent) != 1 {
		t.Fatalf("Since() returned %d entries, want 1", len(recent))
	}
	if recent[0].Event.Kind != state.EventDisappeared {
		t.Errorf("unexpected event kind %v", recent[0].Event.Kind)
	}
}

func TestRecord_EmptySlice_NoOp(t *testing.T) {
	h := history.New(5)
	h.Record(nil)
	if h.Len() != 0 {
		t.Errorf("expected 0 entries after recording nil, got %d", h.Len())
	}
}
