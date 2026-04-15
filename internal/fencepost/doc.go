// Package fencepost provides a lightweight checkpoint tracker for named
// scan intervals.
//
// A Post records the wall-clock time each time a named event is marked and
// exposes whether the elapsed gap since the last mark exceeds a configured
// threshold. This is useful for detecting stalled or missed scan cycles in
// the portwatch daemon without introducing heavy scheduling machinery.
//
// Example:
//
//	fp := fencepost.New(30 * time.Second)
//	fp.Mark("port-scan")
//	// … later …
//	if fp.Overdue("port-scan") {
//		log.Println("scan interval overdue")
//	}
package fencepost
