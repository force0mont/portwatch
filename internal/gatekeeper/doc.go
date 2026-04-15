// Package gatekeeper provides a combined admit/deny guard for alert delivery.
//
// It enforces a per-key rate limit over a rolling time window and applies an
// automatic cooldown period whenever a key is denied, preventing bursts of
// suppressed events from immediately retrying once the window resets.
//
// Typical usage:
//
//	g := gatekeeper.New(gatekeeper.Config{
//		MaxPerWindow:      5,
//		Window:            time.Minute,
//		CooldownAfterDeny: 30 * time.Second,
//	})
//
//	if d := g.Admit("tcp:8080"); d.Admitted {
//		// deliver alert
//	}
package gatekeeper
