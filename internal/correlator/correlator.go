package correlator

import (
	"sync"
	"time"

	"github.com/iamcalledrob/portwatch/internal/scanner"
)

// Correlator groups related port events that occur within a short time
// window into a single correlation group identified by a string key.
type Correlator struct {
	mu      sync.Mutex
	window  time.Duration
	groups  map[string]*Group
	clock   func() time.Time
}

// Group holds a set of ports that appeared together within the correlation window.
type Group struct {
	Key       string
	Ports     []scanner.Port
	FirstSeen time.Time
	LastSeen  time.Time
}

// New returns a Correlator with the given time window.
func New(window time.Duration) *Correlator {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock func() time.Time) *Correlator {
	return &Correlator{
		window: window,
		groups: make(map[string]*Group),
		clock:  clock,
	}
}

// Record associates a port with a correlation key. If an active group
// for that key exists within the window it is updated; otherwise a new
// group is started. The current group is returned.
func (c *Correlator) Record(key string, p scanner.Port) *Group {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.clock()
	g, ok := c.groups[key]
	if !ok || now.Sub(g.LastSeen) > c.window {
		g = &Group{
			Key:       key,
			FirstSeen: now,
		}
		c.groups[key] = g
	}
	g.Ports = append(g.Ports, p)
	g.LastSeen = now
	return g
}

// Flush returns all groups whose window has expired and removes them
// from internal state.
func (c *Correlator) Flush() []*Group {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.clock()
	var expired []*Group
	for key, g := range c.groups {
		if now.Sub(g.LastSeen) > c.window {
			expired = append(expired, g)
			delete(c.groups, key)
		}
	}
	return expired
}

// Len returns the number of active correlation groups.
func (c *Correlator) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.groups)
}
