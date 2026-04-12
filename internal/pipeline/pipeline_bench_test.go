package pipeline_test

import (
	"bytes"
	"context"
	"testing"
	"time"
)

// BenchmarkPipeline_Run measures the overhead of a full tick cycle with an
// empty /proc tree so that scanner.Scan returns quickly.
func BenchmarkPipeline_Run(b *testing.B) {
	var buf bytes.Buffer
	p := buildPipeline(b, &buf)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// prime one tick before measuring
	runFor(ctx, p, 5*time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runFor(ctx, p, 5*time.Millisecond)
	}
}

func buildPipeline(tb testing.TB, buf *bytes.Buffer) interface{ Run(context.Context) error } {
	tb.Helper()
	return buildPipelineInternal(tb, buf)
}

func runFor(ctx context.Context, p interface{ Run(context.Context) error }, d time.Duration) {
	ctx2, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	p.Run(ctx2) //nolint:errcheck
}
