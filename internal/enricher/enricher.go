// Package enricher attaches contextual metadata to scanner port entries
// before they are evaluated by the rules engine or emitted as alerts.
package enricher

import (
	"fmt"
	"net"

	"github.com/user/portwatch/internal/scanner"
)

// Entry is a port entry decorated with additional context.
type Entry struct {
	scanner.Port
	// Hostname is the reverse-DNS name for the listening address, if resolvable.
	Hostname string `json:"hostname,omitempty"`
	// ServiceName is the IANA service name for the port number, if known.
	ServiceName string `json:"service_name,omitempty"`
	// Label is a human-readable string combining protocol, port and service.
	Label string `json:"label"`
}

// Enricher decorates raw scanner ports with metadata.
type Enricher struct {
	lookupAddr func(string) ([]string, error)
	lookupPort func(network, service string) (int, error)
}

// New returns an Enricher that uses the standard library DNS resolver.
func New() *Enricher {
	return &Enricher{
		lookupAddr: net.LookupAddr,
		lookupPort: net.LookupPort,
	}
}

// newWithResolvers constructs an Enricher with injected resolver functions
// for deterministic unit testing.
func newWithResolvers(
	lookupAddr func(string) ([]string, error),
	lookupPort func(string, string) (int, error),
) *Enricher {
	return &Enricher{lookupAddr: lookupAddr, lookupPort: lookupPort}
}

// Enrich takes a scanner.Port and returns an Entry with resolved metadata.
// DNS and service-name failures are silently ignored; the entry is still
// returned with whatever fields could be populated.
func (e *Enricher) Enrich(p scanner.Port) Entry {
	ent := Entry{Port: p}

	// Reverse-DNS lookup.
	if names, err := e.lookupAddr(p.IP); err == nil && len(names) > 0 {
		ent.Hostname = names[0]
	}

	// IANA service name lookup.
	proto := string(p.Protocol)
	if name, err := lookupServiceName(p.Port, proto); err == nil {
		ent.ServiceName = name
	}

	// Build a stable human-readable label.
	if ent.ServiceName != "" {
		ent.Label = fmt.Sprintf("%s/%d (%s)", proto, p.Port, ent.ServiceName)
	} else {
		ent.Label = fmt.Sprintf("%s/%d", proto, p.Port)
	}

	return ent
}

// lookupServiceName maps well-known port numbers to IANA service names using
// a small built-in table so the enricher works without external lookups.
func lookupServiceName(port uint16, proto string) (string, error) {
	key := fmt.Sprintf("%s/%d", proto, port)
	if name, ok := wellKnown[key]; ok {
		return name, nil
	}
	return "", fmt.Errorf("unknown")
}
