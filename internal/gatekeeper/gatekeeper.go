// Package gatekeeper combines rate-limiting, circuit-breaking, and cooldown
// logic into a single admit/deny decision for outbound alert delivery.
package gatekeeper

import (
	"sync"
	"time"
)

// Decision describes why a port event was admitted or denied.
type Decision struct {
	Admitted bool
	Reason   string
}

// Config holds tuneable parameters for the Gatekeeper.
type Config struct {
	// MaxPerWindow is the maximum number of admits allowed per window.
	MaxPerWindow int
	// Window is the rolling time window for the rate limit.
	Window time.Duration
	// CooldownAfterDeny is how long a key stays suppressed after being denied.
	CooldownAfterDeny time.Duration
}

type entry struct {
	count     int
	windowEnd time.Time
	deniedAt  time.Time
}

// Gatekeeper admits or denies alert delivery for a given key.
type Gatekeeper struct {
	mu     sync.Mutex
	cfg    Config
	state  map[string]*entry
	nowFn  func() time.Time
}

// New returns a Gatekeeper with the given config using real wall-clock time.
func New(cfg Config) *Gatekeeper {
	return newWithClock(cfg, time.Now)
}

func newWithClock(cfg Config, nowFn func() time.Time) *Gatekeeper {
	if cfg.MaxPerWindow <= 0 {
		cfg.MaxPerWindow = 10
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.CooldownAfterDeny <= 0 {
		cfg.CooldownAfterDeny = 30 * time.Second
	}
	return &Gatekeeper{
		cfg:   cfg,
		state: make(map[string]*entry),
		nowFn: nowFn,
	}
}

// Admit returns a Decision for the given key (e.g. "tcp:8080").
func (g *Gatekeeper) Admit(key string) Decision {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := g.nowFn()
	e, ok := g.state[key]
	if !ok {
		e = &entry{}
		g.state[key] = e
	}

	// Still within cooldown from a previous denial.
	if !e.deniedAt.IsZero() && now.Before(e.deniedAt.Add(g.cfg.CooldownAfterDeny)) {
		return Decision{Admitted: false, Reason: "cooldown"}
	}

	// Reset window if expired.
	if now.After(e.windowEnd) {
		e.count = 0
		e.windowEnd = now.Add(g.cfg.Window)
	}

	if e.count >= g.cfg.MaxPerWindow {
		e.deniedAt = now
		return Decision{Admitted: false, Reason: "rate_limit"}
	}

	e.count++
	return Decision{Admitted: true, Reason: "ok"}
}

// Reset clears all state for the given key.
func (g *Gatekeeper) Reset(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.state, key)
}
