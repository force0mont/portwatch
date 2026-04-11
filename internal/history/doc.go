// Package history provides a thread-safe, capacity-bounded ring buffer of
// port-scan [state.Event] entries.
//
// The [History] type is intended to be embedded in the watcher loop so that
// operators can query recent port-change activity without needing to persist
// events to an external store.
//
// # Usage
//
//	h := history.New(200)          // keep last 200 events
//	h.Record(events)               // called each scan cycle
//	recent := h.Since(time.Now().Add(-5 * time.Minute))
//
// All methods are safe for concurrent use.
package history
