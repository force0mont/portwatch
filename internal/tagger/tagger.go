// Package tagger assigns human-readable labels to scanner ports
// based on well-known port numbers and configurable overrides.
package tagger

import "sync"

// wellKnown maps common port numbers to service names.
var wellKnown = map[uint16]string{
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Tagger labels ports with service names.
type Tagger struct {
	mu        sync.RWMutex
	overrides map[uint16]string
}

// New returns a Tagger with no custom overrides.
func New() *Tagger {
	return &Tagger{
		overrides: make(map[uint16]string),
	}
}

// Override registers a custom label for the given port, taking
// precedence over the built-in well-known table.
func (t *Tagger) Override(port uint16, label string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.overrides[port] = label
}

// Tag returns the label for port. It checks overrides first,
// then the well-known table, and falls back to "unknown".
func (t *Tagger) Tag(port uint16) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if label, ok := t.overrides[port]; ok {
		return label
	}
	if label, ok := wellKnown[port]; ok {
		return label
	}
	return "unknown"
}

// Known reports whether port has a recognised label (override or
// well-known). Ports tagged as "unknown" return false.
func (t *Tagger) Known(port uint16) bool {
	return t.Tag(port) != "unknown"
}
