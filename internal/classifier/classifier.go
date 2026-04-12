// Package classifier categorises a scanned port into a risk tier based on
// its protocol, port number, and address family.
package classifier

import "github.com/jwhittle933/portwatch/internal/scanner"

// Tier represents a risk classification level.
type Tier int

const (
	// TierLow indicates a well-known, expected service on a standard port.
	TierLow Tier = iota
	// TierMedium indicates an uncommon or registered port that warrants attention.
	TierMedium
	// TierHigh indicates an ephemeral, unknown, or suspicious listener.
	TierHigh
)

// String returns a human-readable label for the tier.
func (t Tier) String() string {
	switch t {
	case TierLow:
		return "low"
	case TierMedium:
		return "medium"
	case TierHigh:
		return "high"
	default:
		return "unknown"
	}
}

// Classifier assigns a risk Tier to a scanned port.
type Classifier struct {
	ephemeralMin uint16
	ephemeralMax uint16
}

// New returns a Classifier using the standard Linux ephemeral port range
// (32768–60999).
func New() *Classifier {
	return &Classifier{ephemeralMin: 32768, ephemeralMax: 60999}
}

// NewWithEphemeralRange returns a Classifier with a custom ephemeral range.
func NewWithEphemeralRange(min, max uint16) *Classifier {
	return &Classifier{ephemeralMin: min, ephemeralMax: max}
}

// Classify returns the risk Tier for the given port entry.
func (c *Classifier) Classify(p scanner.Port) Tier {
	port := p.Port

	// Ephemeral ports are always high risk when listening.
	if port >= c.ephemeralMin && port <= c.ephemeralMax {
		return TierHigh
	}

	// Well-known ports (0-1023) on loopback are low risk.
	if port < 1024 && p.Address == "127.0.0.1" {
		return TierLow
	}

	// Well-known ports on a public interface are medium risk.
	if port < 1024 {
		return TierMedium
	}

	// Registered ports (1024-32767) are medium risk.
	return TierMedium
}
