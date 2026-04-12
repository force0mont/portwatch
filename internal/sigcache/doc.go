// Package sigcache implements a time-bounded signature cache for port-set
// snapshots.
//
// During each scan cycle the pipeline computes a deterministic hash (signature)
// of the current port list. sigcache stores the most-recent signature and
// reports whether an incoming signature matches, allowing the watcher to skip
// expensive downstream stages — rule evaluation, alerting, notification — when
// the port landscape is unchanged.
//
// Entries expire after a configurable TTL so that a long-lived identical
// snapshot does not suppress periodic health reporting indefinitely.
//
// Usage:
//
//	cache := sigcache.New(30 * time.Second)
//
//	if cache.Match(currentSig) {
//	    return // nothing changed
//	}
//	cache.Set(currentSig)
//	// ... run downstream pipeline stages
package sigcache
