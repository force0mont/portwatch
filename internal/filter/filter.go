// Package filter provides port filtering utilities for portwatch.
// It allows callers to suppress or whitelist specific ports, CIDRs,
// or process names before events reach the alerter.
package filter

import (
	"fmt"
	"net"
)

// Rule describes a single suppression rule.
type Rule struct {
	// Port, if non-zero, matches only this port number.
	Port uint16 `json:"port,omitempty"`
	// CIDR, if non-empty, matches addresses within the given network.
	CIDR string `json:"cidr,omitempty"`
	// Protocol restricts the rule to "tcp" or "udp". Empty matches both.
	Protocol string `json:"protocol,omitempty"`

	parsedNet *net.IPNet
}

// Filter holds a compiled set of suppression rules.
type Filter struct {
	rules []Rule
}

// New compiles a slice of Rules into a Filter.
// Returns an error if any CIDR cannot be parsed.
func New(rules []Rule) (*Filter, error) {
	compiled := make([]Rule, len(rules))
	for i, r := range rules {
		if r.CIDR != "" {
			_, ipNet, err := net.ParseCIDR(r.CIDR)
			if err != nil {
				return nil, fmt.Errorf("filter: invalid CIDR %q: %w", r.CIDR, err)
			}
			r.parsedNet = ipNet
		}
		compiled[i] = r
	}
	return &Filter{rules: compiled}, nil
}

// Suppressed returns true when the given port/protocol/addr combination
// is matched by at least one rule in the filter.
func (f *Filter) Suppressed(port uint16, proto, addr string) bool {
	ip := net.ParseIP(addr)
	for _, r := range f.rules {
		if r.Port != 0 && r.Port != port {
			continue
		}
		if r.Protocol != "" && r.Protocol != proto {
			continue
		}
		if r.parsedNet != nil && (ip == nil || !r.parsedNet.Contains(ip)) {
			continue
		}
		return true
	}
	return false
}

// Len returns the number of compiled rules.
func (f *Filter) Len() int { return len(f.rules) }
