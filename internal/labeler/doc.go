// Package labeler attaches human-readable labels to scanned port entries
// based on a configurable set of port/protocol rules.
//
// Labels are short strings (e.g. "http", "ssh", "custom-rpc") that downstream
// components — such as the alerter or reporter — can include in output to make
// events easier to understand at a glance.
//
// # Usage
//
//	l, err := labeler.New([]labeler.Rule{
//		{Port: 22,  Protocol: "tcp", Label: "ssh"},
//		{Port: 443, Protocol: "tcp", Label: "https"},
//	})
//
//	label := l.Label(portEntry) // returns "ssh", "https", or ""
//
// Rules can be added or removed at runtime; all operations are safe for
// concurrent use.
package labeler
