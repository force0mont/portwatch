package metrics_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/metrics"
)

func TestNew_ZeroCounters(t *testing.T) {
	m := metrics.New()
	snap := m.Snapshot()

	if snap.ScansTotal != 0 {
		t.Errorf("expected ScansTotal=0, got %d", snap.ScansTotal)
	}
	if snap.AlertsTotal != 0 {
		t.Errorf("expected AlertsTotal=0, got %d", snap.AlertsTotal)
	}
	if snap.PortsSeen != 0 {
		t.Errorf("expected PortsSeen=0, got %d", snap.PortsSeen)
	}
}

func TestIncScans(t *testing.T) {
	m := metrics.New()
	m.IncScans()
	m.IncScans()

	if got := m.Snapshot().ScansTotal; got != 2 {
		t.Errorf("expected ScansTotal=2, got %d", got)
	}
}

func TestIncAlerts(t *testing.T) {
	m := metrics.New()
	m.IncAlerts()

	if got := m.Snapshot().AlertsTotal; got != 1 {
		t.Errorf("expected AlertsTotal=1, got %d", got)
	}
}

func TestAddPorts(t *testing.T) {
	m := metrics.New()
	m.AddPorts(5)
	m.AddPorts(3)

	if got := m.Snapshot().PortsSeen; got != 8 {
		t.Errorf("expected PortsSeen=8, got %d", got)
	}
}

func TestSnapshot_CollectedAt(t *testing.T) {
	before := time.Now()
	m := metrics.New()
	snap := m.Snapshot()
	after := time.Now()

	if snap.CollectedAt.Before(before) || snap.CollectedAt.After(after) {
		t.Errorf("CollectedAt %v not between %v and %v", snap.CollectedAt, before, after)
	}
}

func TestSnapshot_UptimeNonEmpty(t *testing.T) {
	m := metrics.New()
	snap := m.Snapshot()

	if snap.Uptime == "" {
		t.Error("expected non-empty Uptime string")
	}
}
