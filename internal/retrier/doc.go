// Package retrier provides a context-aware retry helper for portwatch
// operations that may transiently fail, such as webhook delivery or
// file-based baseline persistence.
//
// # Usage
//
//	r := retrier.New(3, 500*time.Millisecond)
//	err := r.Do(ctx, func() error {
//		return sendWebhook(payload)
//	})
//
// Do calls the supplied function up to the configured number of times,
// sleeping between attempts. It returns immediately on success or when
// the context is cancelled, and returns the last error if all attempts
// are exhausted.
package retrier
