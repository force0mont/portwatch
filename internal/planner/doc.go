// Package planner provides scan-scheduling logic for the portwatch daemon.
//
// A Planner wraps a base scan interval with optional jitter to spread load
// and tracks whether individual scan cycles were completed on time. It does
// not drive timers itself; callers use Next to obtain a wait duration and
// Mark to record when a scan actually fired.
//
// Typical usage:
//
//	p := planner.New(30*time.Second, 2*time.Second)
//	for {
//		select {
//		case t := <-time.After(p.Next()):
//			p.Mark(t)
//			// run scan …
//		}
//	}
package planner
