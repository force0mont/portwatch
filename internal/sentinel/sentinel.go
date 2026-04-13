// Package sentinel watches for a specific set of "always-alert" ports that
// should never appear as listeners under any circumstances. Any match bypasses
// the normal rule engine and is immediately flagged as critical.
package sentinel

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Entry describes a port/protocol pair that is unconditionally forbidden.
type Entry struct {
	Port     uint16
	Protocol string // "tcp" or "udp"
}

// Match is returned when a scanned port triggers a sentinel rule.
type Match struct {
	Entry
	Addr string
}

func (m Match) String() string {
	return fmt.Sprintf("sentinel hit: %s/%d on %s", m.Protocol, m.Port, m.Addr)
}

// Sentinel holds the set of forbidden port/protocol pairs.
type Sentinel struct {
	mu      sync.RWMutex
	entries map[Entry]struct{}
}

// New returns a Sentinel pre-loaded with the given entries.
func New(entries []Entry) *Sentinel {
	s := &Sentinel{entries: make(map[Entry]struct{}, len(entries))}
	for _, e := range entries {
		s.entries[e] = struct{}{}
	}
	return s
}

// Add registers an additional forbidden entry at runtime.
func (s *Sentinel) Add(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[e] = struct{}{}
}

// Remove unregisters a forbidden entry.
func (s *Sentinel) Remove(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, e)
}

// Check evaluates a slice of scanned ports and returns every Match found.
// Ports that do not appear in the sentinel set are ignored.
func (s *Sentinel) Check(ports []scanner.Port) []Match {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matches []Match
	for _, p := range ports {
		key := Entry{Port: p.Port, Protocol: p.Protocol}
		if _, ok := s.entries[key]; ok {
			matches = append(matches, Match{Entry: key, Addr: p.Addr})
		}
	}
	return matches
}

// Len returns the number of registered forbidden entries.
func (s *Sentinel) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
