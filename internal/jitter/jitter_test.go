package jitter

import (
	"testing"
	"time"
)

// deterministicSource always returns the same value so tests are repeatable.
type deterministicSource struct{ val int64 }

func (d *deterministicSource) Int63n(n int64) int64 {
	if d.val >= n {
		return n - 1
	}
	return d.val
}

func TestApply_ZeroFactor_ReturnsBase(t *testing.T) {
	j := newWithSource(&deterministicSource{val: 50}, 0)
	base := 100 * time.Millisecond
	if got := j.Apply(base); got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApply_ZeroBase_ReturnsZero(t *testing.T) {
	j := newWithSource(&deterministicSource{val: 50}, 0.2)
	if got := j.Apply(0); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestApply_WithinBounds(t *testing.T) {
	const factor = 0.2
	base := 1000 * time.Millisecond
	low := time.Duration(float64(base) * (1 - factor))
	high := time.Duration(float64(base) * (1 + factor))

	// Run with many synthetic source values to cover the full range.
	for v := int64(0); v < 400; v++ {
		j := newWithSource(&deterministicSource{val: v}, factor)
		got := j.Apply(base)
		if got < low || got > high {
			t.Fatalf("Apply(%v) = %v, want in [%v, %v] (src=%d)", base, got, low, high, v)
		}
	}
}

func TestApply_FactorClampedAboveOne(t *testing.T) {
	j := newWithSource(&deterministicSource{val: 0}, 1.5)
	if j.Factor() >= 1.0 {
		t.Fatalf("factor should be clamped below 1, got %v", j.Factor())
	}
}

func TestApply_FactorClampedBelowZero(t *testing.T) {
	j := newWithSource(&deterministicSource{val: 0}, -0.5)
	if j.Factor() != 0 {
		t.Fatalf("factor should be clamped to 0, got %v", j.Factor())
	}
}

func TestNew_NotNil(t *testing.T) {
	if j := New(0.1); j == nil {
		t.Fatal("New returned nil")
	}
}
