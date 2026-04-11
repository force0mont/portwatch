// Package reporter formats and writes periodic scan summaries.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Format controls how reports are serialised.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Reporter writes periodic metric summaries to an io.Writer.
type Reporter struct {
	w      io.Writer
	format Format
}

// New returns a Reporter that writes to os.Stdout in the given format.
func New(format Format) *Reporter {
	return NewWithWriter(os.Stdout, format)
}

// NewWithWriter returns a Reporter that writes to w.
func NewWithWriter(w io.Writer, format Format) *Reporter {
	return &Reporter{w: w, format: format}
}

// Summary holds the data emitted in each periodic report.
type Summary struct {
	CollectedAt  time.Time `json:"collected_at"`
	TotalScans   int64     `json:"total_scans"`
	TotalAlerts  int64     `json:"total_alerts"`
	CurrentPorts int       `json:"current_ports"`
}

// Report builds a Summary from snap and writes it to the configured writer.
func (r *Reporter) Report(snap metrics.Snapshot) error {
	s := Summary{
		CollectedAt:  snap.CollectedAt,
		TotalScans:   snap.Scans,
		TotalAlerts:  snap.Alerts,
		CurrentPorts: snap.Ports,
	}

	switch r.format {
	case FormatJSON:
		return r.writeJSON(s)
	default:
		return r.writeText(s)
	}
}

func (r *Reporter) writeJSON(s Summary) error {
	enc := json.NewEncoder(r.w)
	return enc.Encode(s)
}

func (r *Reporter) writeText(s Summary) error {
	_, err := fmt.Fprintf(
		r.w,
		"[%s] scans=%d alerts=%d ports=%d\n",
		s.CollectedAt.Format(time.RFC3339),
		s.TotalScans,
		s.TotalAlerts,
		s.CurrentPorts,
	)
	return err
}
