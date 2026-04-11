package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/reporter"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func fixedSnap() metrics.Snapshot {
	return metrics.Snapshot{
		CollectedAt: fixedTime,
		Scans:       42,
		Alerts:      3,
		Ports:       7,
	}
}

func TestReport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewWithWriter(&buf, reporter.FormatText)

	if err := r.Report(fixedSnap()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"scans=42", "alerts=3", "ports=7"} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q missing %q", out, want)
		}
	}
}

func TestReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewWithWriter(&buf, reporter.FormatJSON)

	if err := r.Report(fixedSnap()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var s reporter.Summary
	if err := json.NewDecoder(&buf).Decode(&s); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if s.TotalScans != 42 {
		t.Errorf("want TotalScans=42, got %d", s.TotalScans)
	}
	if s.TotalAlerts != 3 {
		t.Errorf("want TotalAlerts=3, got %d", s.TotalAlerts)
	}
	if s.CurrentPorts != 7 {
		t.Errorf("want CurrentPorts=7, got %d", s.CurrentPorts)
	}
	if !s.CollectedAt.Equal(fixedTime) {
		t.Errorf("want CollectedAt=%v, got %v", fixedTime, s.CollectedAt)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	r := reporter.New(reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
