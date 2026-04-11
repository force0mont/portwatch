// Package audit provides a structured, append-only audit trail for portwatch.
//
// Every significant action performed by the daemon — alerts emitted, ports
// that appeared or disappeared, baseline modifications, and suppressed events
// — is recorded as a newline-delimited JSON entry in the configured audit log
// file.
//
// # Usage
//
//	l, err := audit.New("/var/log/portwatch/audit.log")
//	if err != nil { ... }
//	l.Record(audit.ActionAlertSent, 8080, "tcp", "0.0.0.0", "unexpected listener")
//
// Each entry includes a UTC timestamp, the action type, the port, protocol,
// address, and an optional free-text note.
//
// The Logger is safe for concurrent use.
package audit
