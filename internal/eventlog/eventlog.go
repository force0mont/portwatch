// Package eventlog provides a bounded, time-stamped log of port events
// that can be queried by severity or time range.
package eventlog

import (
	"sync"
	"time"
)

// Level represents the severity of a logged event.
type Level int

const (
	LevelInfo  Level = iota
	LevelAlert Level = iota
)

// Entry is a single record in the event log.
type Entry struct {
	At       time.Time
	Level    Level
	Protocol string
	Port     uint16
	Addr     string
	Message  string
}

// Log is a bounded, thread-safe event log.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	cap     int
	now     func() time.Time
}

// New returns a Log that retains at most capacity entries.
// Oldest entries are evicted when the cap is exceeded.
func New(capacity int) *Log {
	if capacity <= 0 {
		panic("eventlog: capacity must be > 0")
	}
	return newWithClock(capacity, time.Now)
}

func newWithClock(capacity int, now func() time.Time) *Log {
	return &Log{cap: capacity, now: now, entries: make([]Entry, 0, capacity)}
}

// Record appends an entry to the log, evicting the oldest if at capacity.
func (l *Log) Record(level Level, protocol string, port uint16, addr, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := Entry{
		At:       l.now(),
		Level:    level,
		Protocol: protocol,
		Port:     port,
		Addr:     addr,
		Message:  message,
	}
	if len(l.entries) >= l.cap {
		l.entries = append(l.entries[1:], e)
	} else {
		l.entries = append(l.entries, e)
	}
}

// Since returns all entries recorded at or after t.
func (l *Log) Since(t time.Time) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Entry
	for _, e := range l.entries {
		if !e.At.Before(t) {
			out = append(out, e)
		}
	}
	return out
}

// ByLevel returns all entries whose level equals the given level.
func (l *Log) ByLevel(level Level) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Entry
	for _, e := range l.entries {
		if e.Level == level {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the current number of stored entries.
func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}
