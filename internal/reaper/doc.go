// Package reaper provides a TTL-based eviction mechanism for tracking
// which port entries are still actively observed during a portwatch scan
// cycle.
//
// # Overview
//
// A [Reaper] maintains a map of string keys (typically "proto:port") to
// [Entry] values. Callers invoke Touch on each key observed during a scan.
// After the scan completes, calling Reap removes any entry whose LastSeen
// timestamp predates the configured TTL and returns the evicted keys so
// upstream components (e.g. the ledger or state tracker) can react.
//
// # Usage
//
//	r := reaper.New(5 * time.Minute)
//	r.Touch("tcp:8080")
//	evicted := r.Reap() // keys not seen within TTL
package reaper
