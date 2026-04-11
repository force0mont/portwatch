// Package audit provides a structured audit trail for portwatch events,
// recording every action taken (alert sent, baseline updated, suppression
// applied) with a timestamp and contextual metadata.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Action describes what kind of audit event occurred.
type Action string

const (
	ActionAlertSent      Action = "alert_sent"
	ActionSuppressed     Action = "suppressed"
	ActionBaselineUpdate Action = "baseline_update"
	ActionPortAppeared   Action = "port_appeared"
	ActionPortGone       Action = "port_gone"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    Action    `json:"action"`
	Port      uint16    `json:"port"`
	Protocol  string    `json:"protocol"`
	Address   string    `json:"address"`
	Note      string    `json:"note,omitempty"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	mu  sync.Mutex
	out io.Writer
	now func() time.Time
}

// New returns a Logger that writes to the given path, creating it if needed.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return newWithWriter(f), nil
}

// newWithWriter is used in tests to inject an arbitrary writer.
func newWithWriter(w io.Writer) *Logger {
	return &Logger{out: w, now: time.Now}
}

// Record appends an audit entry for the given action.
func (l *Logger) Record(action Action, port uint16, protocol, address, note string) error {
	entry := Entry{
		Timestamp: l.now().UTC(),
		Action:    action,
		Port:      port,
		Protocol:  protocol,
		Address:   address,
		Note:      note,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = fmt.Fprintf(l.out, "%s\n", b)
	return err
}
