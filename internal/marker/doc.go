// Package marker provides an acknowledgement store for open ports.
//
// Operators can explicitly acknowledge known-good listeners so that
// portwatch suppresses further alerts for those ports. Each
// acknowledgement is keyed by (port, protocol) and may carry an
// optional TTL after which it is treated as expired.
//
// Usage:
//
//	m := marker.New()
//	m.Ack(8080, "tcp", "alice", 24*time.Hour, time.Now())
//
//	if m.IsAcked(8080, "tcp", time.Now()) {
//		// suppress alert
//	}
package marker
