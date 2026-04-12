// Package correlator groups port events that occur within a configurable
// time window under a shared correlation key.
//
// This is useful for detecting burst-open patterns — for example, multiple
// ports appearing within a few seconds that may belong to the same process
// or deployment event.
//
// Usage:
//
//	c := correlator.New(3 * time.Second)
//	group := c.Record("deploy-123", port)
//	// later...
//	expired := c.Flush() // returns groups whose window has closed
package correlator
