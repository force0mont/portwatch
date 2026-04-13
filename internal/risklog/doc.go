// Package risklog maintains a running record of open ports together with
// their computed risk scores.
//
// Entries are keyed by (protocol, address) pair so that a port which
// disappears and reappears is treated as the same logical listener.
// Scores can be updated at any time; the original FirstSeen timestamp is
// always preserved.
//
// The TopN method returns entries sorted by descending score, making it
// straightforward to surface the most suspicious listeners for reporting
// or alerting purposes.
//
// All methods are safe for concurrent use.
package risklog
