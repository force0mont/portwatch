// Package alerter handles alert output for unexpected port listeners.
package alerter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/rules"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelAlert Level = "ALERT"
)

// Event represents a single alertable port event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Proto     string
	Port      uint16
	PID       int
	Action    rules.Action
}

// Alerter writes alert events to a configured output.
type Alerter struct {
	out io.Writer
}

// New creates an Alerter that writes to stdout.
func New() *Alerter {
	return &Alerter{out: os.Stdout}
}

// NewWithWriter creates an Alerter that writes to the given writer.
func NewWithWriter(w io.Writer) *Alerter {
	return &Alerter{out: w}
}

// Emit writes a formatted alert event to the output.
func (a *Alerter) Emit(e Event) {
	timestamp := e.Timestamp.Format(time.RFC3339)
	fmt.Fprintf(
		a.out,
		"[%s] %s proto=%s port=%d pid=%d action=%s\n",
		timestamp,
		e.Level,
		e.Proto,
		e.Port,
		e.PID,
		e.Action,
	)
}

// EmitAlert is a convenience method for emitting an ALERT-level event.
func (a *Alerter) EmitAlert(proto string, port uint16, pid int) {
	a.Emit(Event{
		Timestamp: time.Now(),
		Level:     LevelAlert,
		Proto:     proto,
		Port:      port,
		PID:       pid,
		Action:    rules.ActionAlert,
	})
}

// EmitInfo is a convenience method for emitting an INFO-level event.
func (a *Alerter) EmitInfo(proto string, port uint16, pid int) {
	a.Emit(Event{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Proto:     proto,
		Port:      port,
		PID:       pid,
		Action:    rules.ActionAllow,
	})
}

// EmitForAction emits an event at the appropriate level for the given action.
// ALERT-level actions produce an ALERT event; all others produce an INFO event.
func (a *Alerter) EmitForAction(proto string, port uint16, pid int, action rules.Action) {
	level := LevelInfo
	if action == rules.ActionAlert {
		level = LevelAlert
	}
	a.Emit(Event{
		Timestamp: time.Now(),
		Level:     level,
		Proto:     proto,
		Port:      port,
		PID:       pid,
		Action:    action,
	})
}
