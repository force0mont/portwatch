// Package eventlog implements a bounded, time-stamped log of port-change
// events produced by the portwatch scanner pipeline.
//
// Entries are stored in insertion order up to a configurable capacity;
// once full the oldest entry is evicted to make room for the newest.
// All operations are safe for concurrent use.
//
// Typical usage:
//
//	log := eventlog.New(500)
//	log.Record(eventlog.LevelAlert, "tcp", 4444, "0.0.0.0", "unexpected listener")
//
//	// retrieve only alert-level entries
//	alerts := log.ByLevel(eventlog.LevelAlert)
//
//	// retrieve entries from the last 5 minutes
//	recent := log.Since(time.Now().Add(-5 * time.Minute))
package eventlog
