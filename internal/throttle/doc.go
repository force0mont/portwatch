// Package throttle implements per-key event throttling for portwatch.
//
// A Throttle limits how many alert events are forwarded for a given
// port+protocol combination within a configurable sliding window. This
// prevents notification storms when a port flaps rapidly or a rule
// matches a high-frequency transient listener.
//
// Usage:
//
//	th := throttle.New(5, time.Minute)
//	if th.Allow(port, proto) {
//	    // forward the alert
//	}
//
// The zero value is not usable; always construct via New.
package throttle
