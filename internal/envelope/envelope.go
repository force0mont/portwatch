// Package envelope wraps a scanner port entry together with derived
// metadata so that downstream pipeline stages share a single value
// without repeatedly recomputing the same fields.
package envelope

import (
	"fmt"
	"time"

	"github.com/your-org/portwatch/internal/scanner"
)

// Envelope carries a scanned port alongside enrichment data produced by
// earlier pipeline stages.  All fields after Port are optional; a zero
// value means the stage that normally fills it was not wired into the
// pipeline.
type Envelope struct {
	// Port is the raw entry returned by the scanner.
	Port scanner.Port

	// ObservedAt is the wall-clock time at which the port was first seen
	// during the current scan cycle.
	ObservedAt time.Time

	// ServiceName is the human-readable name resolved from the port number
	// (e.g. "http", "ssh").  Empty when resolution was skipped or failed.
	ServiceName string

	// Label is an operator-supplied tag attached by the labeler stage.
	Label string

	// RiskScore is a normalised 0–100 value assigned by the scorecard stage.
	RiskScore int

	// Fingerprint is a short content-hash that identifies the port+protocol
	// combination across scan cycles.
	Fingerprint string

	// Anomalous is true when the anomaly-detection stage considers this port
	// unexpected given the learned baseline.
	Anomalous bool

	// Meta holds arbitrary key/value pairs added by stages that do not
	// warrant a dedicated typed field.
	Meta map[string]string
}

// New returns an Envelope seeded with port and the provided observation
// time.  All enrichment fields are left at their zero values.
func New(p scanner.Port, at time.Time) *Envelope {
	return &Envelope{
		Port:       p,
		ObservedAt: at,
		Meta:       make(map[string]string),
	}
}

// SetMeta stores a key/value pair in the envelope's metadata map.
// It is safe to call on a nil-Meta envelope; the map is initialised
// lazily.
func (e *Envelope) SetMeta(key, value string) {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[key] = value
}

// GetMeta returns the metadata value for key and a boolean indicating
// whether the key was present.
func (e *Envelope) GetMeta(key string) (string, bool) {
	v, ok := e.Meta[key]
	return v, ok
}

// String returns a compact human-readable representation of the envelope
// suitable for log lines and debug output.
func (e *Envelope) String() string {
	return fmt.Sprintf(
		"Envelope{port=%d proto=%s addr=%s risk=%d anomalous=%v label=%q svc=%q fp=%s}",
		e.Port.Port,
		e.Port.Protocol,
		e.Port.Address,
		e.RiskScore,
		e.Anomalous,
		e.Label,
		e.ServiceName,
		e.Fingerprint,
	)
}
