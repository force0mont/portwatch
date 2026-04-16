// Package taproom provides a bounded registry for tracking actively monitored
// port/protocol pairs within portwatch.
//
// A Taproom enforces a configurable maximum number of concurrent taps, making
// it suitable for use as a guard against runaway port registration. Each entry
// is keyed by its protocol and port number, so tcp:80 and udp:80 are treated
// as independent entries.
//
// Typical usage:
//
//	tr := taproom.New(256)
//	if err := tr.Add(8080, "tcp"); err != nil {
//	    // handle ErrTapFull or ErrAlreadyTapped
//	}
//	defer tr.Remove(8080, "tcp")
package taproom
