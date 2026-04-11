// Package suppress provides a time-based suppression mechanism to avoid
// repeatedly alerting on the same port/protocol combination within a
// configurable cooldown window.
package suppress

import (
	"fmt"
	"sync"
	"time"
)

// Clock allows injecting a fake time source in tests.
type Clock func() time.Time

// Suppressor tracks recently alerted keys and suppresses duplicates
// until the cooldown period has elapsed.
type Suppressor struct {
	mu       sync.Mutex
	cooldown time.Duration
	clock    Clock
	last     map[string]time.Time
}

// New returns a Suppressor with the given cooldown duration using the
// real wall clock.
func New(cooldown time.Duration) *Suppressor {
	return newWithClock(cooldown, time.Now)
}

func newWithClock(cooldown time.Duration, clock Clock) *Suppressor {
	return &Suppressor{
		cooldown: cooldown,
		clock:    clock,
		last:     make(map[string]time.Time),
	}
}

// IsSuppressed returns true if an alert for the given port and protocol
// was already emitted within the cooldown window.
func (s *Suppressor) IsSuppressed(port uint16, proto string) bool {
	key := fmt.Sprintf("%s:%d", proto, port)
	now := s.clock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.last[key]; ok && now.Sub(t) < s.cooldown {
		return true
	}

	s.last[key] = now
	return false
}

// Reset clears the suppression record for the given port and protocol,
// allowing the next alert to pass through immediately.
func (s *Suppressor) Reset(port uint16, proto string) {
	key := fmt.Sprintf("%s:%d", proto, port)

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.last, key)
}

// Len returns the number of currently tracked suppression entries.
func (s *Suppressor) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.last)
}
