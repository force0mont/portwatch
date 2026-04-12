// Package decay provides an exponential-decay score tracker for port risk
// accumulation.
//
// Each port key accumulates a floating-point score via Add. The score
// decreases exponentially over time according to a configurable half-life:
// after one half-life the score is halved, after two it is quartered, and so
// on. This allows short-lived noise to fade away while persistent anomalies
// maintain a high score.
//
// Typical usage:
//
//	tr := decay.New(30 * time.Minute)
//	tr.Add("tcp:8080", 5.0)  // port seen — bump score
//	score := tr.Score("tcp:8080") // read current decayed score
//
// Tracker is safe for concurrent use.
package decay
