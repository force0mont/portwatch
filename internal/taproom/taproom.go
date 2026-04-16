// Package taproom provides a port-tap registry that tracks which ports
// have been "tapped" (actively watched) and enforces a maximum tap count.
package taproom

import (
	"errors"
	"fmt"
	"sync"
)

// ErrTapFull is returned when the registry is at capacity.
var ErrTapFull = errors.New("taproom: registry at capacity")

// ErrAlreadyTapped is returned when the port/protocol pair is already registered.
var ErrAlreadyTapped = errors.New("taproom: port already tapped")

// Entry holds metadata for a tapped port.
type Entry struct {
	Port     uint16
	Protocol string
}

// Taproom holds the set of actively watched port/protocol pairs.
type Taproom struct {
	mu      sync.RWMutex
	entries map[string]Entry
	max     int
}

// New creates a Taproom with the given maximum capacity.
// Panics if max is zero.
func New(max int) *Taproom {
	if max <= 0 {
		panic("taproom: max must be greater than zero")
	}
	return &Taproom{
		entries: make(map[string]Entry, max),
		max:     max,
	}
}

func key(port uint16, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Add registers a port/protocol pair. Returns ErrTapFull or ErrAlreadyTapped on failure.
func (t *Taproom) Add(port uint16, proto string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key(port, proto)
	if _, ok := t.entries[k]; ok {
		return ErrAlreadyTapped
	}
	if len(t.entries) >= t.max {
		return ErrTapFull
	}
	t.entries[k] = Entry{Port: port, Protocol: proto}
	return nil
}

// Remove unregisters a port/protocol pair. Returns false if not found.
func (t *Taproom) Remove(port uint16, proto string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key(port, proto)
	if _, ok := t.entries[k]; !ok {
		return false
	}
	delete(t.entries, k)
	return true
}

// Contains reports whether the port/protocol pair is registered.
func (t *Taproom) Contains(port uint16, proto string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.entries[key(port, proto)]
	return ok
}

// Len returns the current number of tapped ports.
func (t *Taproom) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}

// All returns a snapshot of all current entries.
func (t *Taproom) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}
