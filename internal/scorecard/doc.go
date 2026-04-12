// Package scorecard assigns a numeric risk score to each observed port.
//
// Scores are computed from configurable weights applied to signals such as
// whether the service name is recognised, whether the port number falls in the
// ephemeral range (>= 49152), and whether the listener is bound only to the
// loopback interface.
//
// Results are cached in memory so repeated lookups for the same port are O(1).
// Call Evict when a port disappears to release the cached entry.
//
// Usage:
//
//	sc := scorecard.New(scorecard.DefaultWeights())
//	entry := sc.Score(p)
//	fmt.Printf("port %d risk score: %.2f\n", p.Port, entry.Score)
package scorecard
