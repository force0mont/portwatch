// Package notify provides pluggable notification backends for portwatch.
//
// Usage:
//
//	// Write alerts to stderr (default)
//	n := notify.NewStderr()
//	n.Send(notify.Message{
//		Level:    notify.LevelAlert,
//		Port:     8080,
//		Protocol: "tcp",
//		Addr:     "0.0.0.0",
//		Timestamp: time.Now(),
//	})
//
//	// POST JSON to a webhook endpoint
//	wh := notify.NewWebhook("https://hooks.example.com/portwatch")
//	wh.Send(msg)
//
// Implementations satisfy the Notifier interface, so callers can swap
// backends without changing business logic.
package notify
