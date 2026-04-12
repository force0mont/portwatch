// Package sigcache provides a short-lived cache that stores the cryptographic
// signature of a port-set snapshot, allowing the pipeline to skip downstream
// processing when the observed port landscape has not changed between scans.
package sigcache

import (
	"sync"
	"time"
)

// Cache holds the most-recently recorded signature and its timestamp.
type Cache struct {
	mu      sync.Mutex
	sig     string
	recorded time.Time
	ttl     time.Duration
	now     func() time.Time
}

// New returns a Cache whose stored signature expires after ttl.
// A zero ttl means signatures never expire.
func New(ttl time.Duration) *Cache {
	return newWithClock(ttl, time.Now)
}

func newWithClock(ttl time.Duration, now func() time.Time) *Cache {
	return &Cache{ttl: ttl, now: now}
}

// Match reports whether sig equals the cached signature and the entry has not
// yet expired. A miss occurs when the cache is empty, the signature differs, or
// the TTL has elapsed.
func (c *Cache) Match(sig string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sig == "" {
		return false
	}
	if c.ttl > 0 && c.now().Sub(c.recorded) >= c.ttl {
		return false
	}
	return c.sig == sig
}

// Set stores sig as the current cached signature, resetting the TTL clock.
func (c *Cache) Set(sig string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sig = sig
	c.recorded = c.now()
}

// Invalidate clears the cached signature so the next Match call always misses.
func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sig = ""
}
