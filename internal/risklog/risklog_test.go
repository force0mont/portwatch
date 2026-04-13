package risklog

import (
	"testing"
	"time"

	"github.com/iamcathal/portwatch/internal/scanner"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedNow() time.Time { return epoch }

func makePort(proto, addr string) scanner.Port {
	return scanner.Port{Protocol: proto, Address: addr}
}

func TestRecord_And_Len(t *testing.T) {
	l := newWithClock(fixedNow)
	l.Record(makePort("tcp", "0.0.0.0:8080"), 0.5)
	l.Record(makePort("udp", "0.0.0.0:53"), 0.9)
	if got := l.Len(); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}

func TestRecord_UpdatesScore_PreservesFirstSeen(t *testing.T) {
	l := newWithClock(fixedNow)
	p := makePort("tcp", "0.0.0.0:9000")
	l.Record(p, 0.3)

	later := epoch.Add(time.Hour)
	l.now = func() time.Time { return later }
	l.Record(p, 0.8)

	entries := l.TopN(0)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Score != 0.8 {
		t.Errorf("expected score 0.8, got %f", entries[0].Score)
	}
	if !entries[0].FirstSeen.Equal(epoch) {
		t.Errorf("FirstSeen should be preserved as epoch, got %v", entries[0].FirstSeen)
	}
}

func TestTopN_SortedByDescendingScore(t *testing.T) {
	l := newWithClock(fixedNow)
	l.Record(makePort("tcp", "0.0.0.0:80"), 0.2)
	l.Record(makePort("tcp", "0.0.0.0:443"), 0.9)
	l.Record(makePort("udp", "0.0.0.0:53"), 0.5)

	top := l.TopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if top[0].Score < top[1].Score {
		t.Errorf("results not sorted descending: %f < %f", top[0].Score, top[1].Score)
	}
	if top[0].Score != 0.9 {
		t.Errorf("expected top score 0.9, got %f", top[0].Score)
	}
}

func TestTopN_ZeroN_ReturnsAll(t *testing.T) {
	l := newWithClock(fixedNow)
	l.Record(makePort("tcp", "0.0.0.0:22"), 1.0)
	l.Record(makePort("tcp", "0.0.0.0:23"), 0.7)

	if got := len(l.TopN(0)); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	l := newWithClock(fixedNow)
	p := makePort("tcp", "0.0.0.0:8888")
	l.Record(p, 0.6)
	l.Remove(p)
	if l.Len() != 0 {
		t.Errorf("expected 0 entries after remove, got %d", l.Len())
	}
}

func TestRecord_FirstSeen_IsSetAtCreation(t *testing.T) {
	l := newWithClock(fixedNow)
	l.Record(makePort("tcp", "0.0.0.0:3000"), 0.4)
	entries := l.TopN(0)
	if !entries[0].FirstSeen.Equal(epoch) {
		t.Errorf("expected FirstSeen %v, got %v", epoch, entries[0].FirstSeen)
	}
}
