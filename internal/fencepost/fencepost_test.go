package fencepost

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestGap_UnseenName_ReturnsFalse(t *testing.T) {
	p := newWithClock(time.Second, fixedClock(epoch))
	_, ok := p.Gap("scan")
	if ok {
		t.Fatal("expected ok=false for unseen name")
	}
}

func TestMark_And_Gap(t *testing.T) {
	now := epoch
	p := newWithClock(time.Second, func() time.Time { return now })
	p.Mark("scan")
	now = epoch.Add(5 * time.Second)
	gap, ok := p.Gap("scan")
	if !ok {
		t.Fatal("expected ok=true after mark")
	}
	if gap != 5*time.Second {
		t.Fatalf("expected 5s gap, got %v", gap)
	}
}

func TestOverdue_BelowThreshold_ReturnsFalse(t *testing.T) {
	now := epoch
	p := newWithClock(10*time.Second, func() time.Time { return now })
	p.Mark("scan")
	now = epoch.Add(5 * time.Second)
	if p.Overdue("scan") {
		t.Fatal("expected not overdue below threshold")
	}
}

func TestOverdue_AboveThreshold_ReturnsTrue(t *testing.T) {
	now := epoch
	p := newWithClock(5*time.Second, func() time.Time { return now })
	p.Mark("scan")
	now = epoch.Add(10 * time.Second)
	if !p.Overdue("scan") {
		t.Fatal("expected overdue above threshold")
	}
}

func TestOverdue_UnseenName_ReturnsFalse(t *testing.T) {
	p := newWithClock(time.Second, fixedClock(epoch))
	if p.Overdue("never-marked") {
		t.Fatal("unseen name should never be overdue")
	}
}

func TestReset_ClearsCheckpoint(t *testing.T) {
	p := newWithClock(time.Second, fixedClock(epoch))
	p.Mark("scan")
	p.Reset("scan")
	_, ok := p.Gap("scan")
	if ok {
		t.Fatal("expected ok=false after reset")
	}
}

func TestMark_MultipleNames_Independent(t *testing.T) {
	now := epoch
	p := newWithClock(time.Second, func() time.Time { return now })
	p.Mark("a")
	now = epoch.Add(3 * time.Second)
	p.Mark("b")
	now = epoch.Add(6 * time.Second)

	gapA, _ := p.Gap("a")
	gapB, _ := p.Gap("b")
	if gapA != 6*time.Second {
		t.Fatalf("expected gapA=6s, got %v", gapA)
	}
	if gapB != 3*time.Second {
		t.Fatalf("expected gapB=3s, got %v", gapB)
	}
}
