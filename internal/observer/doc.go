// Package observer tracks the presence of open ports across repeated scan
// cycles and computes a stability score for each observed listener.
//
// An Observer is updated each scan cycle via Observe, which accepts the
// current set of live ports. Ports seen in the current scan increment their
// SeenCount; ports previously recorded but absent from the current scan
// increment their MissCount.
//
// The StabilityScore method on Entry returns a value in [0.0, 1.0]:
//
//	1.0 — port has been present in every scan
//	0.5 — port was present in half of all scans
//	0.0 — port has never been seen (or has only been missed)
//
// Observer is safe for concurrent use.
package observer
