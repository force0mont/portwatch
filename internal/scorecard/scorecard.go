// Package scorecard tracks a risk score for each observed port based on
// configurable weights, allowing downstream consumers to prioritise alerts.
package scorecard

import (
	"sync"

	"github.com/yourorg/portwatch/internal/scanner"
)

// Weights controls how much each signal contributes to the final score.
type Weights struct {
	// UnknownService is added when the port has no recognised service name.
	UnknownService float64
	// EphemeralPort is added when the port number is >= 49152.
	EphemeralPort float64
	// LoopbackOnly is subtracted when the address is loopback (lower risk).
	LoopbackOnly float64
}

// DefaultWeights returns sensible out-of-the-box weights.
func DefaultWeights() Weights {
	return Weights{
		UnknownService: 0.4,
		EphemeralPort:  0.3,
		LoopbackOnly:   0.2,
	}
}

// Entry holds the computed score for a single port.
type Entry struct {
	Port  scanner.Port
	Score float64
}

// Scorecard computes and caches risk scores for ports.
type Scorecard struct {
	mu      sync.RWMutex
	weights Weights
	cache   map[string]Entry
}

// New returns a Scorecard using the provided weights.
func New(w Weights) *Scorecard {
	return &Scorecard{
		weights: w,
		cache:   make(map[string]Entry),
	}
}

// Score computes (or retrieves from cache) the risk score for p.
func (s *Scorecard) Score(p scanner.Port) Entry {
	k := key(p)

	s.mu.RLock()
	if e, ok := s.cache[k]; ok {
		s.mu.RUnlock()
		return e
	}
	s.mu.RUnlock()

	score := s.compute(p)
	e := Entry{Port: p, Score: score}

	s.mu.Lock()
	s.cache[k] = e
	s.mu.Unlock()

	return e
}

// Evict removes the cached score for p (e.g. after a port disappears).
func (s *Scorecard) Evict(p scanner.Port) {
	s.mu.Lock()
	delete(s.cache, key(p))
	s.mu.Unlock()
}

func (s *Scorecard) compute(p scanner.Port) float64 {
	var score float64

	if p.Service == "" {
		score += s.weights.UnknownService
	}
	if p.Port >= 49152 {
		score += s.weights.EphemeralPort
	}
	if isLoopback(p.Addr) {
		score -= s.weights.LoopbackOnly
	}
	if score < 0 {
		score = 0
	}
	return score
}

func isLoopback(addr string) bool {
	return addr == "127.0.0.1" || addr == "::1"
}

func key(p scanner.Port) string {
	return p.Protocol + ":" + p.Addr + ":" + string(rune(p.Port))
}
