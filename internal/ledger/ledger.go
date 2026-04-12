// Package ledger maintains a running count of port appearances and
// disappearances over time, providing a simple frequency table that
// other components can query to understand port churn.
package ledger

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds the aggregated counts for a single port/protocol pair.
type Entry struct {
	Key        string
	Appeared   int
	Disappeared int
	LastSeen   time.Time
}

// Ledger tracks appearance and disappearance counts per port/protocol key.
type Ledger struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	clock   func() time.Time
}

// New returns a Ledger using the real wall clock.
func New() *Ledger {
	return newWithClock(time.Now)
}

func newWithClock(clock func() time.Time) *Ledger {
	return &Ledger{
		entries: make(map[string]*Entry),
		clock:   clock,
	}
}

// RecordAppeared increments the appeared counter for the given port/protocol.
func (l *Ledger) RecordAppeared(port uint16, proto string) {
	key := fmt.Sprintf("%s:%d", proto, port)
	l.mu.Lock()
	defer l.mu.Unlock()
	e := l.getOrCreate(key)
	e.Appeared++
	e.LastSeen = l.clock()
}

// RecordDisappeared increments the disappeared counter for the given port/protocol.
func (l *Ledger) RecordDisappeared(port uint16, proto string) {
	key := fmt.Sprintf("%s:%d", proto, port)
	l.mu.Lock()
	defer l.mu.Unlock()
	e := l.getOrCreate(key)
	e.Disappeared++
	e.LastSeen = l.clock()
}

// Get returns a copy of the Entry for the given port/protocol, and whether it exists.
func (l *Ledger) Get(port uint16, proto string) (Entry, bool) {
	key := fmt.Sprintf("%s:%d", proto, port)
	l.mu.RLock()
	defer l.mu.RUnlock()
	e, ok := l.entries[key]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all entries.
func (l *Ledger) All() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, *e)
	}
	return out
}

// Reset clears all entries.
func (l *Ledger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = make(map[string]*Entry)
}

func (l *Ledger) getOrCreate(key string) *Entry {
	e, ok := l.entries[key]
	if !ok {
		e = &Entry{Key: key}
		l.entries[key] = e
	}
	return e
}
