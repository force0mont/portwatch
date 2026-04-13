package eventlog

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_PanicsOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	New(0)
}

func TestRecord_And_Len(t *testing.T) {
	l := newWithClock(10, fixedClock(epoch))
	l.Record(LevelInfo, "tcp", 80, "0.0.0.0", "appeared")
	l.Record(LevelAlert, "tcp", 4444, "0.0.0.0", "unexpected")
	if got := l.Len(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestRecord_EvictsOldestWhenFull(t *testing.T) {
	l := newWithClock(2, fixedClock(epoch))
	l.Record(LevelInfo, "tcp", 80, "0.0.0.0", "first")
	l.Record(LevelInfo, "tcp", 443, "0.0.0.0", "second")
	l.Record(LevelAlert, "tcp", 9999, "0.0.0.0", "third")
	if l.Len() != 2 {
		t.Fatalf("expected cap 2, got %d", l.Len())
	}
	all := l.ByLevel(LevelInfo)
	if len(all) != 1 || all[0].Port != 443 {
		t.Fatalf("expected port 443 retained, got %+v", all)
	}
}

func TestSince_FiltersOldEntries(t *testing.T) {
	l := newWithClock(10, fixedClock(epoch))
	l.Record(LevelInfo, "tcp", 80, "0.0.0.0", "old")

	later := epoch.Add(time.Minute)
	clock := later
	l.now = func() time.Time { return clock }
	l.Record(LevelAlert, "tcp", 4444, "0.0.0.0", "new")

	result := l.Since(later)
	if len(result) != 1 || result[0].Port != 4444 {
		t.Fatalf("expected only new entry, got %+v", result)
	}
}

func TestByLevel_ReturnsOnlyMatchingLevel(t *testing.T) {
	l := newWithClock(10, fixedClock(epoch))
	l.Record(LevelInfo, "tcp", 80, "0.0.0.0", "info")
	l.Record(LevelAlert, "tcp", 1234, "0.0.0.0", "alert")
	l.Record(LevelInfo, "udp", 53, "0.0.0.0", "info2")

	alerts := l.ByLevel(LevelAlert)
	if len(alerts) != 1 || alerts[0].Port != 1234 {
		t.Fatalf("unexpected alerts: %+v", alerts)
	}
	infos := l.ByLevel(LevelInfo)
	if len(infos) != 2 {
		t.Fatalf("expected 2 info entries, got %d", len(infos))
	}
}

func TestSince_EmptyResult_WhenAllOld(t *testing.T) {
	l := newWithClock(10, fixedClock(epoch))
	l.Record(LevelInfo, "tcp", 22, "0.0.0.0", "ssh")

	future := epoch.Add(time.Hour)
	result := l.Since(future)
	if len(result) != 0 {
		t.Fatalf("expected empty, got %+v", result)
	}
}
