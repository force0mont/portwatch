// Package risklog records scored port events and exposes a ranked view
// of the highest-risk listeners seen during the current daemon run.
package risklog

import (
	"sort"
	"sync"
	"time"

	"github.com/iamcathal/portwatch/internal/scanner"
)

// Entry holds a single observed port together with the risk score
// assigned to it and the time it was first recorded.
type Entry struct {
	Port      scanner.Port
	Score     float64
	FirstSeen time.Time
}

// Log stores risk-scored port entries and returns them ranked by score.
type Log struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a Log that uses the real wall clock.
func New() *Log {
	return newWithClock(time.Now)
}

func newWithClock(now func() time.Time) *Log {
	return &Log{
		entries: make(map[string]Entry),
		now:     now,
	}
}

// Record upserts an entry for the given port. If the port was already
// recorded the score is updated but FirstSeen is preserved.
func (l *Log) Record(p scanner.Port, score float64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	k := key(p)
	if existing, ok := l.entries[k]; ok {
		existing.Score = score
		l.entries[k] = existing
		return
	}
	l.entries[k] = Entry{Port: p, Score: score, FirstSeen: l.now()}
}

// TopN returns up to n entries sorted by descending score.
func (l *Log) TopN(n int) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Score > out[j].Score
	})
	if n > 0 && n < len(out) {
		return out[:n]
	}
	return out
}

// Remove deletes the entry for the given port if present.
func (l *Log) Remove(p scanner.Port) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key(p))
}

// Len returns the number of distinct ports currently tracked.
func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}

func key(p scanner.Port) string {
	return p.Protocol + ":" + p.Address
}
