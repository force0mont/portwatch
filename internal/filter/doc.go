// Package filter implements suppression rules for portwatch.
//
// A Filter is constructed from a slice of Rule values, each of which may
// specify a port number, an IP network in CIDR notation, and/or a protocol
// ("tcp" or "udp").  A port event is suppressed when it matches every
// non-zero field of at least one rule — i.e. rules are ANDed within a
// single Rule and ORed across the set.
//
// Typical usage:
//
//	f, err := filter.New(cfg.SuppressRules)
//	if err != nil { ... }
//	if f.Suppressed(port, proto, addr) {
//	    continue // skip this event
//	}
//
// The zero-value Filter (no rules) never suppresses anything.
package filter
