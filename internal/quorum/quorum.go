// Package quorum implements a voting mechanism that requires a port to be
// observed across multiple consecutive scans before it is considered stable.
// This prevents transient or short-lived listeners from triggering alerts.
package quorum

import (
	"sync"

	"github.com/yourorg/portwatch/internal/scanner"
)

// Quorum tracks how many consecutive scans each port has appeared in and
// reports whether the port has reached the required confirmation threshold.
type Quorum struct {
	mu        sync.Mutex
	threshold int
	counts    map[string]int
}

// New returns a Quorum that confirms a port after it has been seen in
// threshold consecutive scans. Panics if threshold < 1.
func New(threshold int) *Quorum {
	if threshold < 1 {
		panic("quorum: threshold must be at least 1")
	}
	return &Quorum{
		threshold: threshold,
		counts:    make(map[string]int),
	}
}

// Observe records a single scan observation for the given port.
// It returns true when the port has been seen in at least threshold
// consecutive scans (i.e. it is now confirmed).
func (q *Quorum) Observe(p scanner.Port) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	k := key(p)
	q.counts[k]++
	return q.counts[k] >= q.threshold
}

// Evict removes the tracking state for a port that has disappeared from
// the scan results, resetting its consecutive-seen counter.
func (q *Quorum) Evict(p scanner.Port) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.counts, key(p))
}

// Count returns the current consecutive-seen count for a port.
func (q *Quorum) Count(p scanner.Port) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.counts[key(p)]
}

func key(p scanner.Port) string {
	return p.Protocol + "/" + p.Address
}
