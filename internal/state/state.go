// Package state tracks the set of previously seen open ports so the watcher
// can detect newly appeared or disappeared listeners between scan cycles.
package state

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// PortKey uniquely identifies a listening socket.
type PortKey struct {
	Proto string
	Addr  string
	Port  uint16
	PID   int
}

func keyFrom(p scanner.Port) PortKey {
	return PortKey{Proto: p.Proto, Addr: p.Addr, Port: p.Port, PID: p.PID}
}

// ChangeKind describes whether a port appeared or disappeared.
type ChangeKind string

const (
	Appeared    ChangeKind = "appeared"
	Disappeared ChangeKind = "disappeared"
)

// Change represents a single state transition for a port.
type Change struct {
	Kind ChangeKind
	Port scanner.Port
}

// Tracker holds the last known set of open ports and computes diffs.
type Tracker struct {
	mu   sync.Mutex
	seen map[PortKey]scanner.Port
}

// New returns an initialised Tracker with an empty baseline.
func New() *Tracker {
	return &Tracker{seen: make(map[PortKey]scanner.Port)}
}

// Diff compares current against the stored baseline, returns the list of
// changes, and updates the baseline to current.
func (t *Tracker) Diff(current []scanner.Port) []Change {
	t.mu.Lock()
	defer t.mu.Unlock()

	next := make(map[PortKey]scanner.Port, len(current))
	for _, p := range current {
		next[keyFrom(p)] = p
	}

	var changes []Change

	// Detect newly appeared ports.
	for k, p := range next {
		if _, ok := t.seen[k]; !ok {
			changes = append(changes, Change{Kind: Appeared, Port: p})
		}
	}

	// Detect disappeared ports.
	for k, p := range t.seen {
		if _, ok := next[k]; !ok {
			changes = append(changes, Change{Kind: Disappeared, Port: p})
		}
	}

	t.seen = next
	return changes
}

// Snapshot returns a copy of the current baseline.
func (t *Tracker) Snapshot() []scanner.Port {
	t.mu.Lock()
	defer t.mu.Unlock()
	ports := make([]scanner.Port, 0, len(t.seen))
	for _, p := range t.seen {
		ports = append(ports, p)
	}
	return ports
}
