// Package marker tracks which ports have been explicitly acknowledged
// by an operator, suppressing repeated alerts for known-good listeners.
package marker

import (
	"fmt"
	"sync"
	"time"
)

// Entry records when a port was acknowledged and by whom.
type Entry struct {
	Port      uint16
	Protocol  string
	AckedBy   string
	AckedAt   time.Time
	ExpiresAt time.Time
}

// Marker stores acknowledgement records for open ports.
type Marker struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Marker.
func New() *Marker {
	return &Marker{
		entries: make(map[string]Entry),
	}
}

// Ack records an acknowledgement for the given port/protocol pair.
// ttl controls how long the acknowledgement remains valid; a zero
// ttl means the acknowledgement never expires.
func (m *Marker) Ack(port uint16, protocol, ackedBy string, ttl time.Duration, now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var exp time.Time
	if ttl > 0 {
		exp = now.Add(ttl)
	}
	m.entries[key(port, protocol)] = Entry{
		Port:      port,
		Protocol:  protocol,
		AckedBy:   ackedBy,
		AckedAt:   now,
		ExpiresAt: exp,
	}
}

// IsAcked reports whether port/protocol has a current acknowledgement.
func (m *Marker) IsAcked(port uint16, protocol string, now time.Time) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[key(port, protocol)]
	if !ok {
		return false
	}
	if !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt) {
		return false
	}
	return true
}

// Revoke removes an acknowledgement for port/protocol.
func (m *Marker) Revoke(port uint16, protocol string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, key(port, protocol))
}

// All returns a snapshot of all current (non-expired) entries.
func (m *Marker) All(now time.Time) []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		if !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt) {
			continue
		}
		out = append(out, e)
	}
	return out
}

func key(port uint16, protocol string) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
