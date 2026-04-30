// Package shedder implements load-shedding for port scan events.
// When the system is under pressure, lower-priority events are dropped
// to protect downstream consumers from being overwhelmed.
package shedder

import (
	"sync"
	"time"
)

// Priority classifies how important an event is for shedding decisions.
type Priority int

const (
	PriorityLow    Priority = 0
	PriorityNormal Priority = 1
	PriorityHigh   Priority = 2
)

// Shedder decides whether an incoming event should be accepted or shed
// based on current load and the event's priority.
type Shedder struct {
	mu        sync.Mutex
	clock     func() time.Time
	window    time.Duration
	maxLow    int
	maxNormal int
	counts    map[Priority][]time.Time
}

// New returns a Shedder with the given window and per-priority maximums.
// maxLow and maxNormal are the maximum events allowed per window for each
// priority tier; High-priority events are never shed.
func New(window time.Duration, maxLow, maxNormal int) *Shedder {
	return newWithClock(window, maxLow, maxNormal, time.Now)
}

func newWithClock(window time.Duration, maxLow, maxNormal int, clock func() time.Time) *Shedder {
	return &Shedder{
		clock:     clock,
		window:    window,
		maxLow:    maxLow,
		maxNormal: maxNormal,
		counts:    make(map[Priority][]time.Time),
	}
}

// Allow returns true if the event with the given priority should be
// processed, or false if it should be shed.
func (s *Shedder) Allow(p Priority) bool {
	if p == PriorityHigh {
		return true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.clock()
	s.evict(p, now)
	max := s.maxFor(p)
	if len(s.counts[p]) >= max {
		return false
	}
	s.counts[p] = append(s.counts[p], now)
	return true
}

// Reset clears all counters for all priorities.
func (s *Shedder) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[Priority][]time.Time)
}

func (s *Shedder) evict(p Priority, now time.Time) {
	cutoff := now.Add(-s.window)
	times := s.counts[p]
	i := 0
	for i < len(times) && times[i].Before(cutoff) {
		i++
	}
	s.counts[p] = times[i:]
}

func (s *Shedder) maxFor(p Priority) int {
	switch p {
	case PriorityLow:
		return s.maxLow
	default:
		return s.maxNormal
	}
}
