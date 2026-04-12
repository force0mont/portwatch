package reaper

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) clock {
	return func() time.Time { return t }
}

func TestTouch_CreatesEntry(t *testing.T) {
	r := newWithClock(time.Minute, fixedClock(epoch))
	r.Touch("tcp:8080")
	if r.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.Len())
	}
}

func TestTouch_IncrementsCount(t *testing.T) {
	r := newWithClock(time.Minute, fixedClock(epoch))
	r.Touch("tcp:8080")
	r.Touch("tcp:8080")
	r.mu.Lock()
	count := r.entries["tcp:8080"].Count
	r.mu.Unlock()
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestReap_RetainsRecentEntries(t *testing.T) {
	r := newWithClock(time.Minute, fixedClock(epoch))
	r.Touch("tcp:9090")
	evicted := r.Reap()
	if len(evicted) != 0 {
		t.Fatalf("expected no evictions, got %v", evicted)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", r.Len())
	}
}

func TestReap_EvictsStaleEntries(t *testing.T) {
	now := epoch
	r := newWithClock(time.Minute, func() time.Time { return now })
	r.Touch("tcp:7070")
	now = epoch.Add(2 * time.Minute)
	evicted := r.Reap()
	if len(evicted) != 1 || evicted[0] != "tcp:7070" {
		t.Fatalf("expected [tcp:7070] evicted, got %v", evicted)
	}
	if r.Len() != 0 {
		t.Fatalf("expected 0 entries after reap, got %d", r.Len())
	}
}

func TestReap_MixedFreshAndStale(t *testing.T) {
	now := epoch
	r := newWithClock(time.Minute, func() time.Time { return now })
	r.Touch("tcp:1111")
	now = epoch.Add(30 * time.Second)
	r.Touch("tcp:2222")
	now = epoch.Add(2 * time.Minute)
	evicted := r.Reap()
	if len(evicted) != 1 || evicted[0] != "tcp:1111" {
		t.Fatalf("expected only tcp:1111 evicted, got %v", evicted)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", r.Len())
	}
}

func TestReap_EmptyReaper_ReturnsNil(t *testing.T) {
	r := newWithClock(time.Minute, fixedClock(epoch))
	evicted := r.Reap()
	if len(evicted) != 0 {
		t.Fatalf("expected empty eviction list, got %v", evicted)
	}
}
