// Package buffer implements a generic, thread-safe ring buffer for portwatch
// internal event queues.
//
// A Buffer retains up to N entries; once full, the oldest entry is silently
// evicted to make room for new arrivals. This makes it suitable for bounded
// in-memory event windows where backpressure is undesirable.
//
// Usage:
//
//	b := buffer.New[scanner.Port](256)
//	b.Push("tcp:8080", port)
//	entries := b.All() // oldest-first snapshot
package buffer
