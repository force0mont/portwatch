// Package preselector filters scanner entries before they reach the
// main processing pipeline, dropping ports that match a static
// ignore-list so downstream stages never see well-known noise.
package preselector

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Preselector holds a set of (port, protocol) pairs that should be
// silently dropped before any rule evaluation occurs.
type Preselector struct {
	mu      sync.RWMutex
	ignored map[key]struct{}
}

type key struct {
	port  uint16
	proto string
}

// New returns a Preselector with no ignored entries.
func New() *Preselector {
	return &Preselector{
		ignored: make(map[key]struct{}),
	}
}

// Ignore registers a (port, protocol) pair that should be dropped by
// Filter. Protocol is matched case-insensitively via normalisation at
// insertion time; callers should pass "tcp" or "udp".
func (p *Preselector) Ignore(port uint16, proto string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ignored[key{port: port, proto: proto}] = struct{}{}
}

// Remove unregisters a previously ignored (port, protocol) pair.
// It is a no-op if the pair was never registered.
func (p *Preselector) Remove(port uint16, proto string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.ignored, key{port: port, proto: proto})
}

// Filter returns only those entries from ports that are NOT on the
// ignore-list. The returned slice is always non-nil.
func (p *Preselector) Filter(ports []scanner.Port) []scanner.Port {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]scanner.Port, 0, len(ports))
	for _, port := range ports {
		if _, ignored := p.ignored[key{port: port.Port, proto: port.Protocol}]; !ignored {
			out = append(out, port)
		}
	}
	return out
}

// Len returns the number of currently ignored (port, protocol) pairs.
func (p *Preselector) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.ignored)
}
