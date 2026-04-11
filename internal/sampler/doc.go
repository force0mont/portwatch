// Package sampler implements a deterministic N-of-M sampling gate for
// port scanner entries.
//
// When the same port/protocol pair is observed repeatedly (e.g. on every
// scan tick) it can be noisy to propagate every observation to downstream
// alerting or reporting stages. Sampler lets callers forward only the 1st,
// (N+1)th, (2N+1)th … occurrence, reducing redundant work while still
// guaranteeing that persistent listeners are periodically re-reported.
//
// Usage:
//
//	s := sampler.New(10) // forward every 10th observation
//	for _, e := range entries {
//	    if s.Sample(e) {
//	        downstream.Handle(e)
//	    }
//	}
//
// Sampler is safe for concurrent use.
package sampler
