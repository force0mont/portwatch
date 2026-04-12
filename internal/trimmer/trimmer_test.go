package trimmer

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAdd_BelowCap_AllRetained(t *testing.T) {
	tr := newWithClock(5, 0, fixedClock(epoch))
	for i := 0; i < 3; i++ {
		tr.Add("k", i)
	}
	if got := tr.Len(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestAdd_ExceedsCap_OldestEvicted(t *testing.T) {
	tr := newWithClock(3, 0, fixedClock(epoch))
	for i := 0; i < 5; i++ {
		tr.Add("k", i)
	}
	if got := tr.Len(); got != 3 {
		t.Fatalf("expected 3 after cap eviction, got %d", got)
	}
	// The three newest values should be 2, 3, 4.
	all := tr.All()
	if all[0].Value.(int) != 2 {
		t.Fatalf("expected oldest retained value=2, got %v", all[0].Value)
	}
}

func TestPrune_RemovesExpiredEntries(t *testing.T) {
	now := epoch
	tr := newWithClock(100, 10*time.Second, func() time.Time { return now })

	tr.Add("a", 1)
	tr.Add("b", 2)

	// Advance clock past TTL.
	now = epoch.Add(15 * time.Second)
	tr.Prune()

	if got := tr.Len(); got != 0 {
		t.Fatalf("expected 0 after TTL prune, got %d", got)
	}
}

func TestPrune_KeepsEntriesWithinTTL(t *testing.T) {
	now := epoch
	tr := newWithClock(100, 10*time.Second, func() time.Time { return now })

	tr.Add("a", 1)
	now = epoch.Add(5 * time.Second)
	tr.Add("b", 2)

	// Advance just past the first entry's TTL.
	now = epoch.Add(12 * time.Second)
	tr.Prune()

	if got := tr.Len(); got != 1 {
		t.Fatalf("expected 1 entry remaining, got %d", got)
	}
	if tr.All()[0].Key != "b" {
		t.Fatalf("expected key 'b' to survive, got %q", tr.All()[0].Key)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := newWithClock(10, 0, fixedClock(epoch))
	tr.Add("x", 42)

	snap := tr.All()
	snap[0].Key = "mutated"

	if tr.All()[0].Key != "x" {
		t.Fatal("All() should return an independent copy")
	}
}

func TestNew_ZeroCap_UsesDefault(t *testing.T) {
	tr := New(0, 0)
	if tr.cap != 256 {
		t.Fatalf("expected default cap 256, got %d", tr.cap)
	}
}
