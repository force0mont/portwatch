package alerter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/rules"
)

func fixedEvent(level alerter.Level, proto string, port uint16, pid int, action rules.Action) alerter.Event {
	return alerter.Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     level,
		Proto:     proto,
		Port:      port,
		PID:       pid,
		Action:    action,
	}
}

func TestEmit_AlertEvent(t *testing.T) {
	var buf bytes.Buffer
	a := alerter.NewWithWriter(&buf)

	e := fixedEvent(alerter.LevelAlert, "tcp", 4444, 1234, rules.ActionAlert)
	a.Emit(e)

	output := buf.String()
	if !strings.Contains(output, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", output)
	}
	if !strings.Contains(output, "port=4444") {
		t.Errorf("expected port=4444 in output, got: %s", output)
	}
	if !strings.Contains(output, "pid=1234") {
		t.Errorf("expected pid=1234 in output, got: %s", output)
	}
	if !strings.Contains(output, "proto=tcp") {
		t.Errorf("expected proto=tcp in output, got: %s", output)
	}
}

func TestEmit_InfoEvent(t *testing.T) {
	var buf bytes.Buffer
	a := alerter.NewWithWriter(&buf)

	e := fixedEvent(alerter.LevelInfo, "udp", 53, 99, rules.ActionAllow)
	a.Emit(e)

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected INFO in output, got: %s", output)
	}
	if !strings.Contains(output, "port=53") {
		t.Errorf("expected port=53 in output, got: %s", output)
	}
}

func TestEmitAlert_Convenience(t *testing.T) {
	var buf bytes.Buffer
	a := alerter.NewWithWriter(&buf)
	a.EmitAlert("tcp", 8080, 555)

	output := buf.String()
	if !strings.Contains(output, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", output)
	}
	if !strings.Contains(output, "port=8080") {
		t.Errorf("expected port=8080 in output, got: %s", output)
	}
}

func TestEmitInfo_Convenience(t *testing.T) {
	var buf bytes.Buffer
	a := alerter.NewWithWriter(&buf)
	a.EmitInfo("udp", 123, 42)

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected INFO in output, got: %s", output)
	}
}
