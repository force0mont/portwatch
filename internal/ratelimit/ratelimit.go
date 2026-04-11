// Package ratelimit provides a token-bucket rate limiter for suppressing
// repeated alerts for the same port within a configurable time window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-key alert counts and suppresses events that exceed
// the configured burst within the rolling window.
type Limiter struct {
	mu      sync.Mutex
	window  time.Duration
	burst   int
	buckets map[string]*bucket
	now     func() time.Time
}

type bucket struct {
	count     int
	windowEnd time.Time
}

// New creates a Limiter that allows up to burst alerts per key within window.
func New(window time.Duration, burst int) *Limiter {
	return newWithClock(window, burst, time.Now)
}

func newWithClock(window time.Duration, burst int, now func() time.Time) *Limiter {
	return &Limiter{
		window:  window,
		burst:   burst,
		buckets: make(map[string]*bucket),
		now:     now,
	}
}

// Allow returns true if the event identified by key should be allowed through.
// It returns false when the burst limit has been reached for the current window.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.buckets[key]
	if !ok || now.After(b.windowEnd) {
		l.buckets[key] = &bucket{count: 1, windowEnd: now.Add(l.window)}
		return true
	}

	if b.count >= l.burst {
		return false
	}
	b.count++
	return true
}

// Reset clears the state for a specific key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Purge removes all expired buckets to free memory.
func (l *Limiter) Purge() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	for k, b := range l.buckets {
		if now.After(b.windowEnd) {
			delete(l.buckets, k)
		}
	}
}
