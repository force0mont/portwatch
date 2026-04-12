// Package watchdog provides a liveness monitor for the portwatch scan loop.
//
// A Watchdog expects periodic calls to Beat; if no beat is received within
// the configured timeout the Stuck channel receives a signal, allowing the
// caller to restart the scan loop or emit an alert.
//
// Usage:
//
//	wd := watchdog.New(30 * time.Second)
//	go wd.Run(ctx)
//
//	// inside the scan loop:
//	wd.Beat()
//
//	// in a supervisor goroutine:
//	select {
//	case <-wd.Stuck():
//	    // handle stall
//	}
package watchdog
