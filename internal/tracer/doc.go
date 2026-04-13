// Package tracer records how long each open port has been continuously visible
// across successive scanner sweeps.
//
// For every port passed to Observe, the tracer maintains:
//   - FirstSeen: the wall-clock time of the first observation
//   - LastSeen:  the wall-clock time of the most recent observation
//   - Duration:  LastSeen − FirstSeen
//
// When a port disappears from a scan result, call Remove so that the next time
// the port appears it starts a fresh tracking window.
//
// All methods are safe for concurrent use.
package tracer
