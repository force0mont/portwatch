// Package baseline manages a persisted snapshot of known-good listening ports.
//
// On first run portwatch can capture the current set of open ports as a
// baseline so that subsequent scans only alert on ports that appeared *after*
// the baseline was taken.  The baseline is written to a JSON file on disk and
// reloaded automatically on startup, making it survive daemon restarts.
//
// Typical usage:
//
//	b, err := baseline.New("/var/lib/portwatch/baseline.json")
//	if err != nil { ... }
//
//	// Capture current listeners as the known-good set.
//	_ = b.Set(currentEntries)
//
//	// Later, check whether a newly discovered port was already baselined.
//	if b.Contains(entry) {
//		// skip – this port was present at baseline time
//	}
package baseline
