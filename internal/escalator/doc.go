// Package escalator promotes alert severity based on repeated port events
// within a rolling time window.
//
// An Escalator tracks how many times a given key (typically "protocol:port")
// has been recorded within a configurable window. When the count crosses the
// warningAt threshold the caller receives LevelWarning; crossing criticalAt
// yields LevelCritical. Counts are reset automatically once the window
// expires, so a port that goes quiet will return to LevelNone on its next
// appearance.
//
// Typical usage:
//
//	esc := escalator.New(5*time.Minute, 3, 7)
//	level := esc.Record(fmt.Sprintf("%s:%d", port.Protocol, port.Port))
//	if level >= escalator.LevelCritical {
//		// page on-call
//	}
package escalator
