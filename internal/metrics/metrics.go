// Package metrics provides lightweight in-process counters for portwatch
// operational telemetry, exposed as a snapshot for logging or diagnostics.
package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Snapshot holds a point-in-time copy of all counters.
type Snapshot struct {
	ScansTotal    uint64    `json:"scans_total"`
	AlertsTotal   uint64    `json:"alerts_total"`
	PortsSeen     uint64    `json:"ports_seen"`
	Uptime        string    `json:"uptime"`
	CollectedAt   time.Time `json:"collected_at"`
}

// Metrics holds atomic counters updated during daemon operation.
type Metrics struct {
	mu          sync.Mutex
	start       time.Time
	scansTotal  atomic.Uint64
	alertsTotal atomic.Uint64
	portsSeen   atomic.Uint64
}

// New creates a new Metrics instance with the start time set to now.
func New() *Metrics {
	return &Metrics{start: time.Now()}
}

// IncScans increments the total scan counter by one.
func (m *Metrics) IncScans() {
	m.scansTotal.Add(1)
}

// IncAlerts increments the total alert counter by one.
func (m *Metrics) IncAlerts() {
	m.alertsTotal.Add(1)
}

// AddPorts adds n to the cumulative ports-seen counter.
func (m *Metrics) AddPorts(n uint64) {
	m.portsSeen.Add(n)
}

// Snapshot returns an immutable copy of the current counter values.
func (m *Metrics) Snapshot() Snapshot {
	now := time.Now()
	return Snapshot{
		ScansTotal:  m.scansTotal.Load(),
		AlertsTotal: m.alertsTotal.Load(),
		PortsSeen:   m.portsSeen.Load(),
		Uptime:      now.Sub(m.start).Round(time.Second).String(),
		CollectedAt: now,
	}
}
