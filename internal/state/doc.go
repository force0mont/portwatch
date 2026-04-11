// Package state provides a thread-safe port-state tracker for portwatch.
//
// The Tracker maintains a baseline snapshot of open ports observed during
// the previous scan cycle. On each new cycle the caller passes the freshly
// scanned slice of scanner.Port values to Diff, which returns a list of
// Change events — one for every port that has appeared or disappeared since
// the last call — and atomically replaces the baseline with the new snapshot.
//
// Typical usage inside the watcher loop:
//
//	tr := state.New()
//	for {
//		ports, _ := s.Scan()
//		for _, ch := range tr.Diff(ports) {
//			// emit alert or info event based on ch.Kind and ch.Port
//		}
//	}
package state
