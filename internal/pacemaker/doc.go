// Package pacemaker monitors the cadence of port-scan cycles and
// reports when scans fall behind the configured interval.
//
// # Overview
//
// A Pacemaker is created with a threshold duration. Each time a scan
// completes the caller invokes Beat; the pacemaker compares the
// elapsed time since the previous beat against the threshold and
// increments an internal missed counter when the gap is too large.
//
// # Usage
//
//	pm := pacemaker.New(5 * time.Second)
//	// after each scan:
//	if ok := pm.Beat(); !ok {
//		log.Printf("scan fell behind: %d total missed", pm.Missed())
//	}
package pacemaker
