package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func fixedTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestRecord_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := newWithWriter(&buf)
	l.now = fixedTime

	if err := l.Record(ActionAlertSent, 8080, "tcp", "0.0.0.0", "unexpected listener"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if entry.Action != ActionAlertSent {
		t.Errorf("action = %q, want %q", entry.Action, ActionAlertSent)
	}
	if entry.Port != 8080 {
		t.Errorf("port = %d, want 8080", entry.Port)
	}
	if entry.Protocol != "tcp" {
		t.Errorf("protocol = %q, want tcp", entry.Protocol)
	}
	if entry.Note != "unexpected listener" {
		t.Errorf("note = %q", entry.Note)
	}
	if !entry.Timestamp.Equal(fixedTime()) {
		t.Errorf("timestamp = %v", entry.Timestamp)
	}
}

func TestRecord_MultipleEntries_NewlineSeparated(t *testing.T) {
	var buf bytes.Buffer
	l := newWithWriter(&buf)
	l.now = fixedTime

	_ = l.Record(ActionPortAppeared, 22, "tcp", "0.0.0.0", "")
	_ = l.Record(ActionPortGone, 22, "tcp", "0.0.0.0", "")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for i, line := range lines {
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Errorf("line %d unmarshal: %v", i, err)
		}
	}
}

func TestRecord_EmptyNote_OmittedFromJSON(t *testing.T) {
	var buf bytes.Buffer
	l := newWithWriter(&buf)
	l.now = fixedTime

	_ = l.Record(ActionSuppressed, 443, "tcp", "127.0.0.1", "")

	if strings.Contains(buf.String(), "note") {
		t.Error("expected 'note' field to be omitted when empty")
	}
}

func TestRecord_AllActions_Valid(t *testing.T) {
	actions := []Action{
		ActionAlertSent,
		ActionSuppressed,
		ActionBaselineUpdate,
		ActionPortAppeared,
		ActionPortGone,
	}
	for _, a := range actions {
		var buf bytes.Buffer
		l := newWithWriter(&buf)
		l.now = fixedTime
		if err := l.Record(a, 80, "tcp", "0.0.0.0", ""); err != nil {
			t.Errorf("action %q: %v", a, err)
		}
	}
}
