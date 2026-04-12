// Package profiler tracks how long each observed port has been continuously
// open and derives an age-based risk score.
//
// A port's AgeScore rises from 0 (seen for less than the warn threshold) to
// 50 (between warn and critical thresholds) to 100 (at or above the critical
// threshold).  Scores can be combined with scorecard weights to produce a
// composite risk signal.
//
// Usage:
//
//	pro := profiler.New()
//	entry := pro.Observe(port)   // call each scan cycle
//	fmt.Println(entry.AgeScore)  // 0, 50, or 100
//
// When a port disappears, call Remove to free the tracking entry.
package profiler
