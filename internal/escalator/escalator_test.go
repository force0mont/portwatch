package escalator

import (
	"testing"
	"time"
)

var (
	fixedNow = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_BelowWarning_ReturnsNone(t *testing.T) {
	e := newWithClock(time.Minute, 3, 5, fixedClock(fixedNow))
	level := e.Record("tcp:8080")
	if level != LevelNone {
		t.Fatalf("expected LevelNone, got %d", level)
	}
}

func TestRecord_AtWarning_ReturnsWarning(t *testing.T) {
	e := newWithClock(time.Minute, 3, 5, fixedClock(fixedNow))
	var level Level
	for i := 0; i < 3; i++ {
		level = e.Record("tcp:8080")
	}
	if level != LevelWarning {
		t.Fatalf("expected LevelWarning, got %d", level)
	}
}

func TestRecord_AtCritical_ReturnsCritical(t *testing.T) {
	e := newWithClock(time.Minute, 3, 5, fixedClock(fixedNow))
	var level Level
	for i := 0; i < 5; i++ {
		level = e.Record("tcp:8080")
	}
	if level != LevelCritical {
		t.Fatalf("expected LevelCritical, got %d", level)
	}
}

func TestRecord_WindowExpiry_ResetsCount(t *testing.T) {
	now := fixedNow
	clock := func() time.Time { return now }
	e := newWithClock(time.Minute, 2, 4, clock)

	e.Record("tcp:9090")
	e.Record("tcp:9090") // hits warning

	// advance past window
	now = now.Add(2 * time.Minute)
	level := e.Record("tcp:9090") // should reset
	if level != LevelNone {
		t.Fatalf("expected LevelNone after window reset, got %d", level)
	}
}

func TestRecord_IndependentKeys(t *testing.T) {
	e := newWithClock(time.Minute, 2, 4, fixedClock(fixedNow))
	e.Record("tcp:80")
	e.Record("tcp:80")
	level := e.Record("udp:53")
	if level != LevelNone {
		t.Fatalf("keys should be independent; expected LevelNone for udp:53, got %d", level)
	}
}

func TestReset_ClearsKey(t *testing.T) {
	e := newWithClock(time.Minute, 2, 4, fixedClock(fixedNow))
	e.Record("tcp:443")
	e.Record("tcp:443")
	e.Reset("tcp:443")
	if e.Len() != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", e.Len())
	}
	level := e.Record("tcp:443")
	if level != LevelNone {
		t.Fatalf("expected LevelNone after reset, got %d", level)
	}
}

func TestLen_TracksEntries(t *testing.T) {
	e := newWithClock(time.Minute, 2, 4, fixedClock(fixedNow))
	e.Record("tcp:80")
	e.Record("tcp:443")
	if got := e.Len(); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}
