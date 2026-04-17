package interval

import (
	"testing"
	"time"
)

const (
	minD  = 2 * time.Second
	maxD  = 30 * time.Second
	stepD = 5 * time.Second
)

func newAdj(thresh int) *Adjuster {
	return New(minD, maxD, stepD, thresh)
}

func TestNew_StartsAtMax(t *testing.T) {
	a := newAdj(3)
	if got := a.Current(); got != maxD {
		t.Fatalf("expected %v got %v", maxD, got)
	}
}

func TestRecordAlert_BelowThreshold_NoChange(t *testing.T) {
	a := newAdj(3)
	a.RecordAlert()
	a.RecordAlert()
	if got := a.Current(); got != maxD {
		t.Fatalf("expected %v got %v", maxD, got)
	}
}

func TestRecordAlert_AtThreshold_DecreasesInterval(t *testing.T) {
	a := newAdj(3)
	a.RecordAlert()
	a.RecordAlert()
	a.RecordAlert()
	want := maxD - stepD
	if got := a.Current(); got != want {
		t.Fatalf("expected %v got %v", want, got)
	}
}

func TestRecordAlert_FloorsAtMin(t *testing.T) {
	a := newAdj(1)
	for i := 0; i < 20; i++ {
		a.RecordAlert()
	}
	if got := a.Current(); got != minD {
		t.Fatalf("expected min %v got %v", minD, got)
	}
}

func TestRelax_IncreasesInterval(t *testing.T) {
	a := newAdj(1)
	a.RecordAlert() // drops to 25s
	a.Relax()       // back to 30s
	if got := a.Current(); got != maxD {
		t.Fatalf("expected %v got %v", maxD, got)
	}
}

func TestRelax_CapsAtMax(t *testing.T) {
	a := newAdj(3)
	a.Relax()
	a.Relax()
	if got := a.Current(); got != maxD {
		t.Fatalf("expected cap at %v got %v", maxD, got)
	}
}

func TestReset_RestoresMax(t *testing.T) {
	a := newAdj(1)
	a.RecordAlert()
	a.RecordAlert()
	a.Reset()
	if got := a.Current(); got != maxD {
		t.Fatalf("expected %v after reset got %v", maxD, got)
	}
}

func TestNew_PanicsOnInvalidParams(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	New(0, maxD, stepD, 3)
}
