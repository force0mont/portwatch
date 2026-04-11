// Package digest computes and compares fingerprints of port snapshots,
// allowing the watcher to detect whether the set of open ports has
// changed between two consecutive scans without storing full state.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

// Digester maintains the last computed digest and reports whether a new
// snapshot differs from the previous one.
type Digester struct {
	mu   sync.Mutex
	last string
}

// New returns a zero-value Digester ready for use.
func New() *Digester {
	return &Digester{}
}

// Entry represents a single open-port observation used as digest input.
type Entry struct {
	Proto   string
	Address string
	Port    uint16
}

// Changed returns true when the digest of entries differs from the
// digest computed during the previous call. The first call always
// returns true so that an initial snapshot is treated as a change.
func (d *Digester) Changed(entries []Entry) bool {
	next := compute(entries)
	d.mu.Lock()
	defer d.mu.Unlock()
	if next == d.last {
		return false
	}
	d.last = next
	return true
}

// Last returns the most recently stored digest, or an empty string if
// Changed has never been called.
func (d *Digester) Last() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.last
}

// compute deterministically hashes a slice of entries.
func compute(entries []Entry) string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = fmt.Sprintf("%s|%s|%d", e.Proto, e.Address, e.Port)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, k := range keys {
		_, _ = fmt.Fprintln(h, k)
	}
	return hex.EncodeToString(h.Sum(nil))
}
