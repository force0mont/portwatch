// Package ratelimit implements a per-key token-bucket rate limiter used by
// portwatch to suppress repeated alert emissions for the same port within a
// rolling time window.
//
// Usage:
//
//	// Allow up to 2 alerts per port every 30 seconds.
//	 limiter := ratelimit.New(30*time.Second, 2)
//
//	 if limiter.Allow("tcp:8080") {
//	     alerter.EmitAlert(event)
//	 }
//
// Purge should be called periodically (e.g. each scan cycle) to reclaim
// memory from expired buckets.
package ratelimit
