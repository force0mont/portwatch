// Package rollup provides event-grouping logic for portwatch.
//
// When the same alert key (e.g. "tcp:8080") fires repeatedly within a
// configurable time window, rollup suppresses the individual events and
// emits a single human-readable summary once the occurrence count reaches
// the configured threshold.  After the window expires the counter resets
// and the cycle begins again.
//
// Typical usage:
//
//	rl := rollup.New(30*time.Second, 5)
//	if summary, ok := rl.Record(key); ok {
//	    alerter.EmitAlert(ctx, summary)
//	}
package rollup
