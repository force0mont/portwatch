// Package anomaly detects ports that deviate from established baseline behaviour.
// It compares the current set of observed ports against a known-good baseline
// and emits an Anomaly for each port that appears unexpected.
package anomaly

import (
	"fmt"
	"sync"
	"time"

	"github.com/jwhittle933/portwatch/internal/scanner"
)

// Anomaly describes a single unexpected port observation.
type Anomaly struct {
	Port      scanner.Port
	Reason    string
	DetectedAt time.Time
}

// String returns a human-readable representation of the anomaly.
func (a Anomaly) String() string {
	return fmt.Sprintf("anomaly: %s/%d – %s", a.Port.Protocol, a.Port.Port, a.Reason)
}

// Detector checks observed ports against a whitelist of known ports.
type Detector struct {
	mu    sync.RWMutex
	known map[string]struct{} // key: "proto:port"
	now   func() time.Time
}

// New returns a Detector with the provided known-good ports pre-loaded.
func New(known []scanner.Port) *Detector {
	return newWithClock(known, time.Now)
}

func newWithClock(known []scanner.Port, now func() time.Time) *Detector {
	d := &Detector{
		known: make(map[string]struct{}, len(known)),
		now:   now,
	}
	for _, p := range known {
		d.known[key(p)] = struct{}{}
	}
	return d
}

// Add registers a port as known-good so it will no longer be flagged.
func (d *Detector) Add(p scanner.Port) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.known[key(p)] = struct{}{}
}

// Check evaluates a slice of ports and returns anomalies for any that are
// not present in the known-good set.
func (d *Detector) Check(ports []scanner.Port) []Anomaly {
	d.mu.RLock()
	defer d.mu.RUnlock()

	now := d.now()
	var out []Anomaly
	for _, p := range ports {
		if _, ok := d.known[key(p)]; !ok {
			out = append(out, Anomaly{
				Port:       p,
				Reason:     "port not in known-good baseline",
				DetectedAt: now,
			})
		}
	}
	return out
}

func key(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Port)
}
