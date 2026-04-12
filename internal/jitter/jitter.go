// Package jitter adds randomised delay to periodic operations to avoid
// thundering-herd problems when multiple goroutines fire at the same time.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Source is the interface satisfied by rand.Rand and any test double.
type Source interface {
	Int63n(n int64) int64
}

// Jitter computes randomised durations within a configurable factor of a base
// interval.  A factor of 0.2 means the returned duration will be within ±20 %
// of the base interval.
type Jitter struct {
	mu     sync.Mutex
	src    Source
	factor float64
}

// New returns a Jitter that uses the global random source and the supplied
// factor.  factor must be in the range [0, 1); values outside that range are
// clamped silently.
func New(factor float64) *Jitter {
	return newWithSource(rand.New(rand.NewSource(time.Now().UnixNano())), factor) //nolint:gosec
}

func newWithSource(src Source, factor float64) *Jitter {
	if factor < 0 {
		factor = 0
	}
	if factor >= 1 {
		factor = 0.99
	}
	return &Jitter{src: src, factor: factor}
}

// Apply returns a duration in the range
//
//	[base*(1-factor), base*(1+factor)]
//
// The result is always ≥ 0.
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if base <= 0 || j.factor == 0 {
		return base
	}

	delta := int64(float64(base) * j.factor)
	if delta == 0 {
		return base
	}

	j.mu.Lock()
	offset := j.src.Int63n(delta*2) - delta
	j.mu.Unlock()

	result := int64(base) + offset
	if result < 0 {
		return 0
	}
	return time.Duration(result)
}

// Factor returns the configured jitter factor.
func (j *Jitter) Factor() float64 { return j.factor }
