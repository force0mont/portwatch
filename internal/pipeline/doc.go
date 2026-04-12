// Package pipeline provides a composable processing chain that connects the
// scanner, state tracker, filter, deduplication, rate-limiter and alerter
// components into a single cohesive scan-and-alert loop.
//
// Typical usage:
//
//	p, err := pipeline.New(pipeline.Config{
//		Scanner:   s,
//		State:     st,
//		Rules:     r,
//		Filter:    f,
//		Dedupe:    d,
//		RateLimit: rl,
//		Alerter:   a,
//		Interval:  5 * time.Second,
//	})
//	p.Run(ctx)
package pipeline
