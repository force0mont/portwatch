// Package dedupe provides event deduplication for portwatch.
//
// A Dedupe instance tracks (port, protocol, address, kind) tuples and
// suppresses repeated occurrences within a configurable time window.
// This prevents alert storms when the same unexpected listener is detected
// across multiple consecutive scans.
//
// Usage:
//
//	d := dedupe.New(10 * time.Second)
//	if !d.IsDuplicate(port, "alert") {
//	    // emit the alert
//	}
//
// The deduplication window is a sliding window: each new unique event
// resets its own timer independently of other events.
package dedupe
