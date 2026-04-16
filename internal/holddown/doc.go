// Package holddown prevents alert flapping by suppressing re-alerts for a
// port that has recently been seen. A port is "held down" from the moment it
// first appears until it has been absent for at least the configured quiet
// period, at which point the hold is released and the next appearance will
// trigger a fresh alert.
//
// Typical usage:
//
//	h := holddown.New(30 * time.Second)
//
//	// on each scan tick, for every active port:
//	if h.Seen(port, proto) {
//		// first appearance — emit alert
//	}
//
//	// for ports that disappeared this tick:
//	h.Gone(port, proto)
//
//	// periodically release stale holds:
//	h.Prune()
package holddown
