// Package fence implements a sliding-window trip-wire.
//
// A Fence counts events per key within a configurable time window and
// signals when the count reaches a threshold. Once tripped, the caller
// is responsible for deciding how to respond (e.g. escalate, suppress,
// or reset the counter via Reset).
//
// Typical usage:
//
//	f := fence.New(30*time.Second, 5)
//	if trip, ok := f.Record(portKey); ok {
//		log.Printf("port %s tripped fence (%d events)", trip.Key, trip.Count)
//		f.Reset(portKey)
//	}
package fence
