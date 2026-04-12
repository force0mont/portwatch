package profiler

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/scanner"
)

var (
	t0    = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	testPort = scanner.Port{Protocol: "tcp", Address: "0.0.0.0", Port: 8080}
)

func fixedAt(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestObserve_FirstCall_ZeroDuration(t *testing.T) {
	p := newWithClock(fixedAt(t0), 5*time.Minute, time.Hour)
	e := p.Observe(testPort)

	if e.Duration != 0 {
		t.Fatalf("expected zero duration on first observe, got %v", e.Duration)
	}
	if e.AgeScore != 0 {
		t.Fatalf("expected score 0, got %d", e.AgeScore)
	}
}

func TestObserve_AfterWarnThreshold_Score50(t *testing.T) {
	clock := fixedAt(t0)
	p := newWithClock(clock, 5*time.Minute, time.Hour)
	p.Observe(testPort)

	// Advance clock past warn threshold.
	clock = fixedAt(t0.Add(6 * time.Minute))
	p.clock = clock
	e := p.Observe(testPort)

	if e.AgeScore != 50 {
		t.Fatalf("expected score 50, got %d", e.AgeScore)
	}
}

func TestObserve_AfterCriticalThreshold_Score100(t *testing.T) {
	clock := fixedAt(t0)
	p := newWithClock(clock, 5*time.Minute, time.Hour)
	p.Observe(testPort)

	p.clock = fixedAt(t0.Add(2 * time.Hour))
	e := p.Observe(testPort)

	if e.AgeScore != 100 {
		t.Fatalf("expected score 100, got %d", e.AgeScore)
	}
}

func TestRemove_ClearsEntry(t *testing.T) {
	p := newWithClock(fixedAt(t0), 5*time.Minute, time.Hour)
	p.Observe(testPort)
	p.Remove(testPort)

	snap := p.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot after remove, got %d entries", len(snap))
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	p := newWithClock(fixedAt(t0), 5*time.Minute, time.Hour)
	p.Observe(testPort)

	s1 := p.Snapshot()
	s1["mutated"] = Entry{}
	s2 := p.Snapshot()

	if _, ok := s2["mutated"]; ok {
		t.Fatal("snapshot modification leaked into profiler state")
	}
}

func TestObserve_DifferentPorts_Independent(t *testing.T) {
	p := newWithClock(fixedAt(t0), 5*time.Minute, time.Hour)
	other := scanner.Port{Protocol: "udp", Address: "0.0.0.0", Port: 9090}

	p.Observe(testPort)
	p.Observe(other)

	snap := p.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
}
