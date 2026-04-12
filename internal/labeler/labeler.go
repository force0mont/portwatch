// Package labeler attaches human-readable labels to scanned ports based on
// configurable tag rules and well-known service mappings.
package labeler

import (
	"fmt"
	"sync"

	"github.com/patrickward/portwatch/internal/scanner"
)

// Rule associates a label with a port/protocol pair.
type Rule struct {
	Port     uint16
	Protocol string // "tcp" or "udp"
	Label    string
}

// Labeler maps ports to labels.
type Labeler struct {
	mu    sync.RWMutex
	rules map[string]string // key: "proto:port" → label
}

// New returns a Labeler pre-loaded with the given rules.
func New(rules []Rule) (*Labeler, error) {
	l := &Labeler{
		rules: make(map[string]string, len(rules)),
	}
	for _, r := range rules {
		if r.Label == "" {
			return nil, fmt.Errorf("labeler: empty label for port %d/%s", r.Port, r.Protocol)
		}
		if r.Protocol != "tcp" && r.Protocol != "udp" {
			return nil, fmt.Errorf("labeler: unknown protocol %q", r.Protocol)
		}
		l.rules[key(r.Protocol, r.Port)] = r.Label
	}
	return l, nil
}

// Label returns the label for the given port entry, or an empty string when
// no rule matches.
func (l *Labeler) Label(p scanner.PortEntry) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.rules[key(p.Protocol, p.Port)]
}

// Add registers or replaces a label rule at runtime.
func (l *Labeler) Add(r Rule) error {
	if r.Label == "" {
		return fmt.Errorf("labeler: empty label for port %d/%s", r.Port, r.Protocol)
	}
	l.mu.Lock()
	l.rules[key(r.Protocol, r.Port)] = r.Label
	l.mu.Unlock()
	return nil
}

// Remove deletes a label rule. It is a no-op if the rule does not exist.
func (l *Labeler) Remove(protocol string, port uint16) {
	l.mu.Lock()
	delete(l.rules, key(protocol, port))
	l.mu.Unlock()
}

func key(protocol string, port uint16) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
