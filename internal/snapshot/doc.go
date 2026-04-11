// Package snapshot provides an immutable, thread-safe representation of the
// open port set captured at a specific point in time.
//
// A Snapshot is created from a slice of scanner.Port values and exposes
// read-only operations for membership tests, enumeration, and diffing against
// a subsequent snapshot.
//
// Typical usage:
//
//	prev := snapshot.New(firstScan)
//	// … wait for next scan interval …
//	next := snapshot.New(secondScan)
//	appeared, disappeared := prev.Diff(next)
//
// Snapshot is safe for concurrent use by multiple goroutines.
package snapshot
