// Package sampler provides port-scan sampling that reduces noise by
// only forwarding every Nth observation of the same port/protocol pair.
package sampler

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Sampler forwards a port entry only once every N calls for the same key.
type Sampler struct {
	mu     sync.Mutex
	n      int
	counts map[string]int
}

// New returns a Sampler that passes through every Nth observation.
// n must be >= 1; values less than 1 are treated as 1 (pass-through).
func New(n int) *Sampler {
	if n < 1 {
		n = 1
	}
	return &Sampler{
		n:      n,
		counts: make(map[string]int),
	}
}

// Sample returns true when the entry should be forwarded downstream.
// It increments an internal counter per (addr:port/proto) key and
// returns true only when count % n == 1 (i.e. the 1st, (n+1)th, … call).
func (s *Sampler) Sample(e scanner.Entry) bool {
	key := keyFor(e)
	s.mu.Lock()
	s.counts[key]++
	c := s.counts[key]
	s.mu.Unlock()
	return c%s.n == 1
}

// Reset clears all counters, causing the next observation of every key
// to be forwarded regardless of previous history.
func (s *Sampler) Reset() {
	s.mu.Lock()
	s.counts = make(map[string]int)
	s.mu.Unlock()
}

// Len returns the number of distinct keys currently tracked.
func (s *Sampler) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.counts)
}

func keyFor(e scanner.Entry) string {
	return e.Addr + ":" + itoa(e.Port) + "/" + e.Protocol
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
