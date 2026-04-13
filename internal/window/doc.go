// Package window implements a thread-safe sliding-window counter.
//
// It is useful for rate-limiting and frequency analysis: callers add
// occurrences keyed by an arbitrary string and query how many events fell
// within the most recent duration.
//
// Typical usage:
//
//	w := window.New(30 * time.Second)
//
//	// record an event and get the running total
//	count := w.Add("192.168.1.1:8080/tcp")
//	if count > threshold {
//		// take action
//	}
//
// Expired buckets are pruned lazily on every Add or Count call, so there is
// no background goroutine to manage.
package window
