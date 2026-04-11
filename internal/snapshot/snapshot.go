// Package snapshot provides point-in-time capture and comparison of
// open port sets, enabling portwatch to detect changes between scans.
package snapshot

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds a captured port entry alongside the time it was first seen.
type Entry struct {
	Port      scanner.Port
	FirstSeen time.Time
}

// Snapshot is an immutable, thread-safe capture of the open port set at a
// specific point in time.
type Snapshot struct {
	mu        sync.RWMutex
	capturedAt time.Time
	entries    map[string]Entry // key: "proto:addr:port"
}

// New creates a Snapshot from the given ports, recording the current time.
func New(ports []scanner.Port) *Snapshot {
	return newAt(ports, time.Now())
}

func newAt(ports []scanner.Port, t time.Time) *Snapshot {
	s := &Snapshot{
		capturedAt: t,
		entries:    make(map[string]Entry, len(ports)),
	}
	for _, p := range ports {
		k := key(p)
		s.entries[k] = Entry{Port: p, FirstSeen: t}
	}
	return s
}

// CapturedAt returns the time the snapshot was taken.
func (s *Snapshot) CapturedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.capturedAt
}

// Contains reports whether the given port is present in the snapshot.
func (s *Snapshot) Contains(p scanner.Port) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.entries[key(p)]
	return ok
}

// Len returns the number of entries in the snapshot.
func (s *Snapshot) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// All returns a copy of all entries in the snapshot.
func (s *Snapshot) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

// Diff computes the ports that appeared in next but not in s (appeared) and
// ports present in s but missing from next (disappeared).
func (s *Snapshot) Diff(next *Snapshot) (appeared, disappeared []scanner.Port) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	next.mu.RLock()
	defer next.mu.RUnlock()

	for k, e := range next.entries {
		if _, ok := s.entries[k]; !ok {
			appeared = append(appeared, e.Port)
		}
	}
	for k, e := range s.entries {
		if _, ok := next.entries[k]; !ok {
			disappeared = append(disappeared, e.Port)
		}
	}
	return
}

func key(p scanner.Port) string {
	return p.Protocol + ":" + p.Address + ":" + itoa(p.Port)
}

func itoa(n uint16) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	buf := [5]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = digits[n%10]
		n /= 10
	}
	return string(buf[pos:])
}
