package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func buildPipeline(t *testing.T, buf *bytes.Buffer) *pipeline.Pipeline {
	t.Helper()
	s, _ := scanner.NewWithProcPath(t.TempDir())
	st := state.New()
	r, _ := rules.New(nil)
	f, _ := filter.New(nil)
	d := dedupe.New(500 * time.Millisecond)
	rl := ratelimit.New(10, time.Second)
	a := alerter.NewWithWriter(buf)
	return pipeline.New(pipeline.Config{
		Scanner:   s,
		State:     st,
		Rules:     r,
		Filter:    f,
		Dedupe:    d,
		RateLimit: rl,
		Alerter:   a,
		Interval:  20 * time.Millisecond,
	})
}

func TestNew_NotNil(t *testing.T) {
	var buf bytes.Buffer
	p := buildPipeline(t, &buf)
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestRun_CancelsCleanly(t *testing.T) {
	var buf bytes.Buffer
	p := buildPipeline(t, &buf)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	err := p.Run(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_DoesNotBlockIndefinitely(t *testing.T) {
	var buf bytes.Buffer
	p := buildPipeline(t, &buf)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- p.Run(ctx) }()
	time.Sleep(40 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("pipeline did not stop after context cancel")
	}
}
