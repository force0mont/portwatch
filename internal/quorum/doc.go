// Package quorum provides a confirmation gate for newly discovered ports.
//
// A port must be observed in a configurable number of consecutive scan
// cycles before it is considered stable and forwarded to the alert pipeline.
// This eliminates noise caused by ephemeral sockets that appear and vanish
// within a single scan interval.
//
// Usage:
//
//	q := quorum.New(3) // require 3 consecutive sightings
//
//	for _, p := range scanResults {
//	    if q.Observe(p) {
//	        // port has been seen 3 times in a row — emit alert
//	    }
//	}
//
//	// When a port disappears from scan results, reset its counter:
//	for _, gone := range disappeared {
//	    q.Evict(gone)
//	}
package quorum
