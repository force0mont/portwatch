// Package interval provides adaptive scan interval adjustment based on
// recent alert activity. When alerts are frequent the interval shrinks;
// when quiet it relaxes back toward the configured maximum.
package interval

import (
	"sync"
	"time"
)

// Adjuster tracks alert frequency and returns a recommended scan interval.
type Adjuster struct {
	mu      sync.Mutex
	min     time.Duration
	max     time.Duration
	current time.Duration
	step    time.Duration
	alerts  int
	thresh  int
}

// New returns an Adjuster with the given min/max bounds, step size, and alert
// threshold. When the number of recorded alerts reaches thresh the interval
// is decreased by step (floored at min). Each call to Relax increases it by
// step (capped at max).
func New(min, max, step time.Duration, alertThresh int) *Adjuster {
	if min <= 0 || max < min || step <= 0 || alertThresh <= 0 {
		panic("interval: invalid parameters")
	}
	return &Adjuster{
		min:     min,
		max:     max,
		current: max,
		step:    step,
		thresh:  alertThresh,
	}
}

// RecordAlert notes that an alert was emitted. When accumulated alerts reach
// the threshold the interval is tightened and the counter resets.
func (a *Adjuster) RecordAlert() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.alerts++
	if a.alerts >= a.thresh {
		a.alerts = 0
		a.current -= a.step
		if a.current < a.min {
			a.current = a.min
		}
	}
}

// Relax nudges the interval one step toward the maximum. Call after a quiet
// scan cycle.
func (a *Adjuster) Relax() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.current += a.step
	if a.current > a.max {
		a.current = a.max
	}
}

// Current returns the recommended scan interval.
func (a *Adjuster) Current() time.Duration {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.current
}

// Reset restores the interval to max and clears the alert counter.
func (a *Adjuster) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.current = a.max
	a.alerts = 0
}
