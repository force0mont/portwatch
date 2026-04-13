// Package topology groups observed ports by address family and protocol,
// providing a structured view of the current listener landscape.
package topology

import (
	"net"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Group holds ports that share the same protocol and address class.
type Group struct {
	Protocol string
	Class    string // "loopback", "private", "public"
	Ports    []scanner.Port
}

// Topology maps (protocol, class) keys to port groups.
type Topology struct {
	mu     sync.RWMutex
	groups map[string]*Group
}

// New returns an empty Topology.
func New() *Topology {
	return &Topology{groups: make(map[string]*Group)}
}

// Build replaces all groups with a fresh classification of ports.
func (t *Topology) Build(ports []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.groups = make(map[string]*Group)
	for _, p := range ports {
		cls := classify(p.Addr)
		k := p.Protocol + "/" + cls
		g, ok := t.groups[k]
		if !ok {
			g = &Group{Protocol: p.Protocol, Class: cls}
			t.groups[k] = g
		}
		g.Ports = append(g.Ports, p)
	}
}

// Groups returns a snapshot of all current groups.
func (t *Topology) Groups() []Group {
	t.mu.RLock()
	defer t.mu.RUnlock()

	out := make([]Group, 0, len(t.groups))
	for _, g := range t.groups {
		copy := *g
		copy.Ports = append([]scanner.Port(nil), g.Ports...)
		out = append(out, copy)
	}
	return out
}

// Len returns the total number of tracked ports across all groups.
func (t *Topology) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	n := 0
	for _, g := range t.groups {
		n += len(g.Ports)
	}
	return n
}

func classify(addr string) string {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "unknown"
	}
	if ip.IsLoopback() {
		return "loopback"
	}
	if isPrivate(ip) {
		return "private"
	}
	return "public"
}

var privateRanges = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"fc00::/7",
}

func isPrivate(ip net.IP) bool {
	for _, cidr := range privateRanges {
		_, n, err := net.ParseCIDR(cidr)
		if err == nil && n.Contains(ip) {
			return true
		}
	}
	return false
}
